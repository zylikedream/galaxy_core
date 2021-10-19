/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package connection

import (
	"compress/flate"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/AlexStocks/log4go"
	"github.com/golang/snappy"

	"github.com/pkg/errors"
)

var (
	launchTime = time.Now()
	connID     uint32

// ErrInvalidConnection = errors.New("connection has been closed.")
)

// for zip compress
type writeFlusher struct {
	flusher *flate.Writer
	lock    sync.Mutex
}

func (t *writeFlusher) Write(p []byte) (int, error) {
	var (
		n   int
		err error
	)
	t.lock.Lock()
	defer t.lock.Unlock()
	n, err = t.flusher.Write(p)
	if err != nil {
		return n, errors.WithStack(err)
	}
	if err := t.flusher.Flush(); err != nil {
		return 0, errors.WithStack(err)
	}

	return n, nil
}

/////////////////////////////////////////
// getty tcp connection
/////////////////////////////////////////

type GettyTCPConn struct {
	gettyConn
	reader io.Reader
	writer io.Writer
	conn   net.Conn
}

// create gettyTCPConn
func NewGettyTCPConn(conn net.Conn) *GettyTCPConn {
	if conn == nil {
		panic("newGettyTCPConn(conn):@conn is nil")
	}
	var localAddr, peerAddr string
	//  check conn.LocalAddr or conn.RemoetAddr is nil to defeat panic on 2016/09/27
	if conn.LocalAddr() != nil {
		localAddr = conn.LocalAddr().String()
	}
	if conn.RemoteAddr() != nil {
		peerAddr = conn.RemoteAddr().String()
	}

	return &GettyTCPConn{
		conn:   conn,
		reader: io.Reader(conn),
		writer: io.Writer(conn),
		gettyConn: gettyConn{
			id:       atomic.AddUint32(&connID, 1),
			rTimeout: netIOTimeout,
			wTimeout: netIOTimeout,
			local:    localAddr,
			peer:     peerAddr,
			compress: CompressNone,
		},
	}
}

// set compress type(tcp: zip/snappy, websocket:zip)
func (t *GettyTCPConn) SetCompressType(c CompressType) {
	switch c {
	case CompressNone, CompressZip, CompressBestSpeed, CompressBestCompression, CompressHuffman:
		ioReader := io.Reader(t.conn)
		t.reader = flate.NewReader(ioReader)

		ioWriter := io.Writer(t.conn)
		w, err := flate.NewWriter(ioWriter, int(c))
		if err != nil {
			panic(fmt.Sprintf("flate.NewReader(flate.DefaultCompress) = err(%s)", err))
		}
		t.writer = &writeFlusher{flusher: w}

	case CompressSnappy:
		ioReader := io.Reader(t.conn)
		t.reader = snappy.NewReader(ioReader)
		ioWriter := io.Writer(t.conn)
		t.writer = snappy.NewBufferedWriter(ioWriter)

	default:
		panic(fmt.Sprintf("illegal comparess type %d", c))
	}
	t.compress = c
}

// tcp connection read
func (t *GettyTCPConn) Recv(p []byte) (int, error) {
	var (
		err         error
		currentTime time.Time
		length      int
	)

	// set read timeout deadline
	if t.compress == CompressNone && t.rTimeout > 0 {
		// Optimization: update read deadline only if more than 25%
		// of the last read deadline exceeded.
		// See https://github.com/golang/go/issues/15133 for details.
		currentTime = time.Now()
		if currentTime.Sub(t.rLastDeadline) > (t.rTimeout >> 2) {
			if err = t.conn.SetReadDeadline(currentTime.Add(t.rTimeout)); err != nil {
				// just a timeout error
				return 0, errors.WithStack(err)
			}
			t.rLastDeadline = currentTime
		}
	}

	length, err = t.reader.Read(p)
	// log.Debug("now:%s, length:%d, err:%s", currentTime, length, err)
	atomic.AddUint32(&t.readBytes, uint32(length))
	return length, errors.WithStack(err)
	//return length, err
}

// tcp connection write
func (t *GettyTCPConn) Send(pkg interface{}) (int, error) {
	var (
		err         error
		currentTime time.Time
		ok          bool
		p           []byte
		length      int
	)

	if t.compress == CompressNone && t.wTimeout > 0 {
		// Optimization: update write deadline only if more than 25%
		// of the last write deadline exceeded.
		// See https://github.com/golang/go/issues/15133 for details.
		currentTime = time.Now()
		if currentTime.Sub(t.wLastDeadline) > (t.wTimeout >> 2) {
			if err = t.conn.SetWriteDeadline(currentTime.Add(t.wTimeout)); err != nil {
				return 0, errors.WithStack(err)
			}
			t.wLastDeadline = currentTime
		}
	}
	if buffers, ok := pkg.([][]byte); ok {
		netBuf := net.Buffers(buffers)
		if length, err := netBuf.WriteTo(t.conn); err == nil {
			atomic.AddUint32(&t.writeBytes, (uint32)(length))
			atomic.AddUint32(&t.writePkgNum, (uint32)(len(buffers)))
		}
		log.Debug("localAddr: %s, remoteAddr:%s, now:%s, length:%d, err:%s",
			t.conn.LocalAddr(), t.conn.RemoteAddr(), currentTime, length, err)
		return int(length), errors.WithStack(err)
	}

	if p, ok = pkg.([]byte); ok {
		if length, err = t.writer.Write(p); err == nil {
			atomic.AddUint32(&t.writeBytes, (uint32)(len(p)))
			atomic.AddUint32(&t.writePkgNum, 1)
		}
		log.Debug("localAddr: %s, remoteAddr:%s, now:%s, length:%d, err:%s",
			t.conn.LocalAddr(), t.conn.RemoteAddr(), currentTime, length, err)
		return length, errors.WithStack(err)
	}
	return 0, errors.Errorf("illegal @pkg{%#v} type", pkg)
	//return length, err
}

// close tcp connection
func (t *GettyTCPConn) Close(waitSec int) {
	// if tcpConn, ok := t.conn.(*net.TCPConn); ok {
	// tcpConn.SetLinger(0)
	// }

	if t.conn != nil {
		if writer, ok := t.writer.(*snappy.Writer); ok {
			if err := writer.Close(); err != nil {
				log.Error("snappy.Writer.Close() = error{%s}", errors.WithStack(err))
			}
		}
		if conn, ok := t.conn.(*net.TCPConn); ok {
			_ = conn.SetLinger(waitSec)
			_ = conn.Close()
		} else {
			_ = t.conn.(*tls.Conn).Close()

		}
		t.conn = nil
	}
}

func (t *GettyTCPConn) NetConn() net.Conn {
	return t.conn
}

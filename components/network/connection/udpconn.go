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
	"fmt"
	"net"
	"sync/atomic"
	"time"

	log "github.com/AlexStocks/log4go"

	"github.com/pkg/errors"
)

var ErrNullPeerAddr = errors.New("peer address is nil")

type UDPContext struct {
	Pkg      interface{}
	PeerAddr *net.UDPAddr
}

func (c UDPContext) String() string {
	return fmt.Sprintf("{pkg:%#v, peer addr:%s}", c.Pkg, c.PeerAddr)
}

type GettyUDPConn struct {
	gettyConn
	compressType CompressType
	conn         *net.UDPConn // for server
}

// create gettyUDPConn
func NewGettyUDPConn(conn *net.UDPConn) *GettyUDPConn {
	if conn == nil {
		panic("newGettyUDPConn(conn):@conn is nil")
	}

	var localAddr, peerAddr string
	if conn.LocalAddr() != nil {
		localAddr = conn.LocalAddr().String()
	}

	if conn.RemoteAddr() != nil {
		// connected udp
		peerAddr = conn.RemoteAddr().String()
	}

	return &GettyUDPConn{
		conn: conn,
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

func (u *GettyUDPConn) SetCompressType(c CompressType) {
	switch c {
	case CompressNone, CompressZip, CompressBestSpeed, CompressBestCompression, CompressHuffman, CompressSnappy:
		u.compressType = c

	default:
		panic(fmt.Sprintf("illegal comparess type %d", c))
	}
}

// udp connection read
func (u *GettyUDPConn) Recv(p []byte) (int, *net.UDPAddr, error) {
	var (
		err         error
		currentTime time.Time
		length      int
		addr        *net.UDPAddr
	)

	if u.rTimeout > 0 {
		// Optimization: update read deadline only if more than 25%
		// of the last read deadline exceeded.
		// See https://github.com/golang/go/issues/15133 for details.
		currentTime = time.Now()
		if currentTime.Sub(u.rLastDeadline) > (u.rTimeout >> 2) {
			if err = u.conn.SetReadDeadline(currentTime.Add(u.rTimeout)); err != nil {
				return 0, nil, errors.WithStack(err)
			}
			u.rLastDeadline = currentTime
		}
	}

	length, addr, err = u.conn.ReadFromUDP(p) // connected udp also can get return @addr
	log.Debug("ReadFromUDP() = {length:%d, peerAddr:%s, error:%s}", length, addr, err)
	if err == nil {
		atomic.AddUint32(&u.readBytes, uint32(length))
	}

	//return length, addr, err
	return length, addr, errors.WithStack(err)
}

// write udp packet, @ctx should be of type UDPContext
func (u *GettyUDPConn) Send(udpCtx interface{}) (int, error) {
	ctx, ok := udpCtx.(UDPContext)
	if !ok {
		return 0, errors.Errorf("illegal @udpCtx{%s} type, @udpCtx type:%T", udpCtx, udpCtx)
	}
	buf, ok := ctx.Pkg.([]byte)
	if !ok {
		return 0, errors.Errorf("illegal @udpCtx.Pkg{%#v} type", udpCtx)
	}
	peerAddr := ctx.PeerAddr
	if peerAddr == nil {
		return 0, ErrNullPeerAddr
	}

	if u.wTimeout > 0 {
		// Optimization: update write deadline only if more than 25%
		// of the last write deadline exceeded.
		// See https://github.com/golang/go/issues/15133 for details.
		currentTime := time.Now()
		if currentTime.Sub(u.wLastDeadline) > (u.wTimeout >> 2) {
			if err := u.conn.SetWriteDeadline(currentTime.Add(u.wTimeout)); err != nil {
				return 0, errors.WithStack(err)
			}
			u.wLastDeadline = currentTime
		}
	}

	length, _, err := u.conn.WriteMsgUDP(buf, nil, peerAddr)
	if err == nil {
		atomic.AddUint32(&u.writeBytes, (uint32)(len(buf)))
		atomic.AddUint32(&u.writePkgNum, 1)
	}
	log.Debug("WriteMsgUDP(peerAddr:%s) = {length:%d, error:%s}", peerAddr, length, err)

	return length, errors.WithStack(err)
	//return length, err
}

// close udp connection
func (u *GettyUDPConn) Close(_ int) {
	if u.conn != nil {
		u.conn.Close()
		u.conn = nil
	}
}

func (u *GettyUDPConn) NetConn() net.Conn {
	return u.conn
}

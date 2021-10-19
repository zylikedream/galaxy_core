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
	"net"
	"sync/atomic"
	"time"
)

const (
	netIOTimeout = 1e9 // 1s
)

type CompressType int

const (
	CompressNone            CompressType = flate.NoCompression      // 0
	CompressZip                          = flate.DefaultCompression // -1
	CompressBestSpeed                    = flate.BestSpeed          // 1
	CompressBestCompression              = flate.BestCompression    // 9
	CompressHuffman                      = flate.HuffmanOnly        // -2
	CompressSnappy                       = 10
)

type Connection interface {
	ID() uint32
	SetCompressType(CompressType)
	LocalAddr() string
	RemoteAddr() string
	IncReadPkgNum()
	IncWritePkgNum()
	// update session's active time
	UpdateActive()
	// get session's active time
	GetActive() time.Time
	ReadTimeout() time.Duration
	// SetReadTimeout sets deadline for the future read calls.
	SetReadTimeout(time.Duration)
	WriteTimeout() time.Duration
	// SetWriteTimeout sets deadline for the future read calls.
	SetWriteTimeout(time.Duration)
	Send(interface{}) (int, error)
	NetConn() net.Conn
	// don't distinguish between tcp connection and websocket connection. Because
	// gorilla/websocket/conn.go:(Conn)Close also invoke net.Conn.Close
	Close(int)
	// set related session
}

type gettyConn struct {
	id            uint32
	compress      CompressType
	padding1      uint8
	padding2      uint16
	readBytes     uint32        // read bytes
	writeBytes    uint32        // write bytes
	readPkgNum    uint32        // send pkg number
	writePkgNum   uint32        // recv pkg number
	active        int64         // last active, in milliseconds
	rTimeout      time.Duration // network current limiting
	wTimeout      time.Duration
	rLastDeadline time.Time // lastest network read time
	wLastDeadline time.Time // lastest network write time
	local         string    // local address
	peer          string    // peer address
}

func (c *gettyConn) ID() uint32 {
	return c.id
}

func (c *gettyConn) LocalAddr() string {
	return c.local
}

func (c *gettyConn) RemoteAddr() string {
	return c.peer
}

func (c *gettyConn) IncReadPkgNum() {
	atomic.AddUint32(&c.readPkgNum, 1)
}

func (c *gettyConn) IncWritePkgNum() {
	atomic.AddUint32(&c.writePkgNum, 1)
}

func (c *gettyConn) UpdateActive() {
	atomic.StoreInt64(&(c.active), int64(time.Since(launchTime)))
}

func (c *gettyConn) GetActive() time.Time {
	return launchTime.Add(time.Duration(atomic.LoadInt64(&(c.active))))
}

func (c *gettyConn) send(interface{}) (int, error) {
	return 0, nil
}

func (c *gettyConn) close(int) {}

func (c gettyConn) ReadTimeout() time.Duration {
	return c.rTimeout
}

// Pls do not set read deadline for websocket connection. AlexStocks 20180310
// gorilla/websocket/conn.go:NextReader will always fail when got a timeout error.
//
// Pls do not set read deadline when using compression. AlexStocks 20180314.
func (c *gettyConn) SetReadTimeout(rTimeout time.Duration) {
	if rTimeout < 1 {
		panic("@rTimeout < 1")
	}

	c.rTimeout = rTimeout
	if c.wTimeout == 0 {
		c.wTimeout = rTimeout
	}
}

func (c gettyConn) WriteTimeout() time.Duration {
	return c.wTimeout
}

// Pls do not set write deadline for websocket connection. AlexStocks 20180310
// gorilla/websocket/conn.go:NextWriter will always fail when got a timeout error.
//
// Pls do not set write deadline when using compression. AlexStocks 20180314.
func (c *gettyConn) SetWriteTimeout(wTimeout time.Duration) {
	if wTimeout < 1 {
		panic("@wTimeout < 1")
	}

	c.wTimeout = wTimeout
	if c.rTimeout == 0 {
		c.rTimeout = wTimeout
	}
}

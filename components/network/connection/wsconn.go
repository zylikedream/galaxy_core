package connection

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	log "github.com/AlexStocks/log4go"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

/////////////////////////////////////////
// getty websocket connection
/////////////////////////////////////////

type GettyWSConn struct {
	gettyConn
	conn *websocket.Conn
}

// create websocket connection
func NewGettyWSConn(conn *websocket.Conn) *GettyWSConn {
	if conn == nil {
		panic("newGettyWSConn(conn):@conn is nil")
	}
	var localAddr, peerAddr string
	//  check conn.LocalAddr or conn.RemoetAddr is nil to defeat panic on 2016/09/27
	if conn.LocalAddr() != nil {
		localAddr = conn.LocalAddr().String()
	}
	if conn.RemoteAddr() != nil {
		peerAddr = conn.RemoteAddr().String()
	}

	gettyWSConn := &GettyWSConn{
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
	conn.EnableWriteCompression(false)
	conn.SetPingHandler(gettyWSConn.handlePing)
	conn.SetPongHandler(gettyWSConn.handlePong)

	return gettyWSConn
}

// set compress type
func (w *GettyWSConn) SetCompressType(c CompressType) {
	switch c {
	case CompressNone, CompressZip, CompressBestSpeed, CompressBestCompression, CompressHuffman:
		w.conn.EnableWriteCompression(true)
		w.conn.SetCompressionLevel(int(c))

	default:
		panic(fmt.Sprintf("illegal comparess type %d", c))
	}
	w.compress = c
}

func (w *GettyWSConn) handlePing(message string) error {
	err := w.WritePong([]byte(message))
	if err == websocket.ErrCloseSent {
		err = nil
	} else if e, ok := err.(net.Error); ok && e.Temporary() {
		err = nil
	}
	if err == nil {
		w.UpdateActive()
	}

	return errors.WithStack(err)
}

func (w *GettyWSConn) handlePong(string) error {
	w.UpdateActive()
	return nil
}

// websocket connection read
func (w *GettyWSConn) Recv() ([]byte, error) {
	// Pls do not set read deadline when using ReadMessage. AlexStocks 20180310
	// gorilla/websocket/conn.go:NextReader will always fail when got a timeout error.
	_, b, e := w.conn.ReadMessage() // the first return value is message type.
	if e == nil {
		atomic.AddUint32(&w.readBytes, (uint32)(len(b)))
	} else {
		if websocket.IsUnexpectedCloseError(e, websocket.CloseGoingAway) {
			log.Warn("websocket unexpected close error: %v", e)
		}
	}

	return b, errors.WithStack(e)
	//return b, e
}

func (w *GettyWSConn) updateWriteDeadline() error {
	var (
		err         error
		currentTime time.Time
	)

	if w.wTimeout > 0 {
		// Optimization: update write deadline only if more than 25%
		// of the last write deadline exceeded.
		// See https://github.com/golang/go/issues/15133 for details.
		currentTime = time.Now()
		if currentTime.Sub(w.wLastDeadline) > (w.wTimeout >> 2) {
			if err = w.conn.SetWriteDeadline(currentTime.Add(w.wTimeout)); err != nil {
				return errors.WithStack(err)
			}
			w.wLastDeadline = currentTime
		}
	}

	return nil
}

// websocket connection write
func (w *GettyWSConn) Send(pkg interface{}) (int, error) {
	var (
		err error
		ok  bool
		p   []byte
	)

	if p, ok = pkg.([]byte); !ok {
		return 0, errors.Errorf("illegal @pkg{%#v} type", pkg)
	}

	w.updateWriteDeadline()
	if err = w.conn.WriteMessage(websocket.BinaryMessage, p); err == nil {
		atomic.AddUint32(&w.writeBytes, (uint32)(len(p)))
		atomic.AddUint32(&w.writePkgNum, 1)
	}
	return len(p), errors.WithStack(err)
	//return len(p), err
}

func (w *GettyWSConn) WritePing() error {
	w.updateWriteDeadline()
	return errors.WithStack(w.conn.WriteMessage(websocket.PingMessage, []byte{}))
}

func (w *GettyWSConn) WritePong(message []byte) error {
	w.updateWriteDeadline()
	return errors.WithStack(w.conn.WriteMessage(websocket.PongMessage, message))
}

// close websocket connection
func (w *GettyWSConn) Close(waitSec int) {
	w.updateWriteDeadline()
	w.conn.WriteMessage(websocket.CloseMessage, []byte("bye-bye!!!"))
	conn := w.conn.UnderlyingConn()
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetLinger(waitSec)
	} else if wsConn, ok := conn.(*tls.Conn); ok {
		wsConn.CloseWrite()
	}
	w.conn.Close()
}

func (w *GettyWSConn) NetConn() net.Conn {
	return w.conn.UnderlyingConn()
}

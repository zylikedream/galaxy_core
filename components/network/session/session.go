package session

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	gxbytes "github.com/dubbogo/gost/bytes"
	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/components/network"
	"github.com/zylikedream/galaxy/components/network/connection"
	"github.com/zylikedream/galaxy/components/network/peer"

	log "github.com/AlexStocks/log4go"
	gxcontext "github.com/dubbogo/gost/context"
	gxtime "github.com/dubbogo/gost/time"
)

var (
	ErrSessionClosed  = errors.New("session Already Closed")
	ErrSessionBlocked = errors.New("session Full Blocked")
)

const (
	maxReadBufLen    = 4 * 1024
	netIOTimeout     = 1e9      // 1s
	period           = 60 * 1e9 // 1 minute
	pendingDuration  = 3e9
	defaultQLen      = 1024
	maxIovecNum      = 10
	MaxWheelTimeSpan = 900e9 // 900s, 15 minute

	defaultSessionName    = "session"
	defaultTCPSessionName = "tcp-session"
	defaultUDPSessionName = "udp-session"
	defaultWSSessionName  = "ws-session"
	defaultWSSSessionName = "wss-session"
	outputFormat          = "session %s, Read Bytes: %d, Write Bytes: %d, Read Pkgs: %d, Write Pkgs: %d"
)

/////////////////////////////////////////
// session
/////////////////////////////////////////

var (
	wheel              *gxtime.Wheel
	sessionClientKey   = "session-client-owner"
	connectPingPackage = []byte("connect-ping")
)

func init() {
	span := 100e6 // 100ms
	buckets := MaxWheelTimeSpan / span
	wheel = gxtime.NewWheel(time.Duration(span), int(buckets)) // wheel longest span is 15 minute
}

func GetTimeWheel() *gxtime.Wheel {
	return wheel
}

// getty base session
type session struct {
	name string
	peer peer.Peer

	// net read Write
	conn     connection.Connection
	listener network.EventListener

	// codec
	reader network.Reader // @reader should be nil when @conn is a gettyWSConn object.
	writer network.Writer

	// write
	wQ chan interface{}

	// handle logic
	maxMsgLen int32

	// heartbeat
	period time.Duration

	// done
	wait time.Duration
	once *sync.Once
	done chan struct{}

	// attribute
	attrs *gxcontext.ValuesContext

	// goroutines sync
	grNum int32
	// read goroutines done signal
	rDone chan struct{}
	lock  sync.RWMutex
}

func newSession(peer peer.Peer, conn connection.Connection) *session {
	ss := &session{
		name: defaultSessionName,
		peer: peer,

		conn: conn,

		maxMsgLen: maxReadBufLen,

		period: period,

		once:  &sync.Once{},
		done:  make(chan struct{}),
		wait:  pendingDuration,
		attrs: gxcontext.NewValuesContext(context.Background()),
		rDone: make(chan struct{}),
	}

	ss.conn.SetWriteTimeout(netIOTimeout)
	ss.conn.SetReadTimeout(netIOTimeout)

	return ss
}

func (s *session) Reset() {
	*s = session{
		name:   defaultSessionName,
		once:   &sync.Once{},
		done:   make(chan struct{}),
		period: period,
		wait:   pendingDuration,
		attrs:  gxcontext.NewValuesContext(context.Background()),
		rDone:  make(chan struct{}),
	}
}

// func (s *session) SetConn(conn net.Conn) { s.gettyConn = newGettyConn(conn) }
func (s *session) Peer() peer.Peer {
	return s.peer
}

func (s *session) Conn() connection.Connection {
	return s.conn
}

func (s *session) Stat() string {
	// stat todo
	return ""
}

// check whether the session has been closed.
func (s *session) IsClosed() bool {
	select {
	case <-s.done:
		return true

	default:
		return false
	}
}

// set maximum package length of every package in (EventListener)OnMessage(@pkgs)
func (s *session) SetMaxMsgLen(length int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.maxMsgLen = int32(length)
}

// set session name
func (s *session) SetName(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.name = name
}

// set EventListener
func (s *session) SetEventListener(listener network.EventListener) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.listener = listener
}

// set package handler
func (s *session) SetPkgHandler(handler network.ReadWriter) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.reader = handler
	s.writer = handler
}

// set Reader
func (s *session) SetReader(reader network.Reader) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.reader = reader
}

// set Writer
func (s *session) SetWriter(writer network.Writer) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.writer = writer
}

// period is in millisecond. Websocket session will send ping frame automatically every peroid.
func (s *session) SetCronPeriod(period int) {
	if period < 1 {
		panic("@period < 1")
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.period = time.Duration(period) * time.Millisecond
}

// set @session's Write queue size
func (s *session) SetWQLen(writeQLen int) {
	if writeQLen < 1 {
		panic("@writeQLen < 1")
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.wQ = make(chan interface{}, writeQLen)
	log.Debug("%s, [session.SetWQLen] wQ{len:%d, cap:%d}", s.Stat(), len(s.wQ), cap(s.wQ))
}

// set maximum wait time when session got error or got exit signal
func (s *session) SetWaitTime(waitTime time.Duration) {
	if waitTime < 1 {
		panic("@wait < 1")
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.wait = waitTime
}

// set attribute of key @session:key
func (s *session) GetAttribute(key interface{}) interface{} {
	s.lock.RLock()
	if s.attrs == nil {
		s.lock.RUnlock()
		return nil
	}
	ret, flag := s.attrs.Get(key)
	s.lock.RUnlock()

	if !flag {
		return nil
	}

	return ret
}

// get attribute of key @session:key
func (s *session) SetAttribute(key interface{}, value interface{}) {
	s.lock.Lock()
	if s.attrs != nil {
		s.attrs.Set(key, value)
	}
	s.lock.Unlock()
}

// delete attribute of key @session:key
func (s *session) RemoveAttribute(key interface{}) {
	s.lock.Lock()
	if s.attrs != nil {
		s.attrs.Delete(key)
	}
	s.lock.Unlock()
}

func (s *session) sessionToken() string {
	if s.IsClosed() || s.conn == nil {
		return "session-closed"
	}

	return fmt.Sprintf("{%s:%s:%d:%s<->%s}",
		s.name, s.Peer().Type(), s.conn.ID(), s.conn.LocalAddr(), s.conn.RemoteAddr())
}

func (s *session) WritePkg(pkg interface{}, timeout time.Duration) error {
	if pkg == nil {
		return fmt.Errorf("@pkg is nil")
	}
	if s.IsClosed() {
		return ErrSessionClosed
	}

	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			rBuf := make([]byte, size)
			rBuf = rBuf[:runtime.Stack(rBuf, false)]
			log.Error("[session.WritePkg] panic session %s: err=%s\n%s", s.sessionToken(), r, rBuf)
		}
	}()

	if timeout <= 0 {
		pkgBytes, err := s.writer.Write(s, pkg)
		if err != nil {
			log.Warn("%s, [session.WritePkg] session.writer.Write(@pkg:%#v) = error:%+v", s.Stat(), pkg, err)
			return errors.WithStack(err)
		}

		var udpCtxPtr *connection.UDPContext
		if udpCtx, ok := pkg.(connection.UDPContext); ok {
			udpCtxPtr = &udpCtx
		} else if udpCtxP, ok := pkg.(*connection.UDPContext); ok {
			udpCtxPtr = udpCtxP
		}
		if udpCtxPtr != nil {
			udpCtxPtr.Pkg = pkgBytes
			pkg = *udpCtxPtr
		} else {
			pkg = pkgBytes
		}
		_, err = s.conn.Send(pkg)
		if err != nil {
			log.Warn("%s, [session.WritePkg] @s.Connection.Write(pkg:%#v) = err:%+v", s.Stat(), pkg, err)
			return errors.WithStack(err)
		}
		return nil
	}
	select {
	case s.wQ <- pkg:
		break // for possible gen a new pkg

	case <-wheel.After(timeout):
		log.Warn("%s, [session.WritePkg] wQ{len:%d, cap:%d}", s.Stat(), len(s.wQ), cap(s.wQ))
		return ErrSessionBlocked
	}

	return nil
}

// for codecs
func (s *session) WriteBytes(pkg []byte) error {
	if s.IsClosed() {
		return ErrSessionClosed
	}

	// s.conn.SetWriteTimeout(time.Now().Add(s.wTimeout))
	if _, err := s.conn.Send(pkg); err != nil {
		return errors.Wrapf(err, "s.Connection.Write(pkg len:%d)", len(pkg))
	}
	return nil
}

// Write multiple packages at once. so we invoke write sys.call just one time.
func (s *session) WriteBytesArray(pkgs ...[]byte) error {
	if s.IsClosed() {
		return ErrSessionClosed
	}
	// s.conn.SetWriteTimeout(time.Now().Add(s.wTimeout))
	if len(pkgs) == 1 {
		// return s.Connection.Write(pkgs[0])
		return s.WriteBytes(pkgs[0])
	}

	// reduce syscall and memcopy for multiple packages
	// if _, ok := s.conn.(*gettyTCPConn); ok {
	if true {
		if _, err := s.conn.Send(pkgs); err != nil {
			return errors.Wrapf(err, "s.Connection.Write(pkgs num:%d)", len(pkgs))
		}
		return nil
	}
	// }

	var length int
	for i := 0; i < len(pkgs); i++ {
		length += len(pkgs[i])
	}

	// merge the pkgs
	//arr = make([]byte, length)
	arrp := gxbytes.GetBytes(length)
	defer gxbytes.PutBytes(arrp)
	arr := *arrp

	l := 0
	for i := 0; i < len(pkgs); i++ {
		copy(arr[l:], pkgs[i])
		l += len(pkgs[i])
	}

	if err := s.WriteBytes(arr); err != nil {
		return errors.WithStack(err)
	}

	num := len(pkgs) - 1
	for i := 0; i < num; i++ {
		s.conn.IncWritePkgNum()
	}

	return nil
}

// func (s *session) RunEventLoop() {
func (s *session) Run() {
	if s.conn == nil || s.listener == nil || s.writer == nil {
		errStr := fmt.Sprintf("session{name:%s, conn:%#v, listener:%#v, writer:%#v}",
			s.name, s.conn, s.listener, s.writer)
		log.Error(errStr)
		panic(errStr)
	}

	if s.wQ == nil {
		s.wQ = make(chan interface{}, defaultQLen)
	}

	// call session opened
	s.conn.UpdateActive()
	if err := s.listener.OnOpen(s); err != nil {
		log.Error("[OnOpen] session %s, error: %#v", s.Stat(), err)
		s.Close()
		return
	}

	// start read/write gr
	atomic.AddInt32(&(s.grNum), 2)
	go s.handleLoop()
	go s.handlePackage()
}

func (s *session) handleLoop() {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			rBuf := make([]byte, size)
			rBuf = rBuf[:runtime.Stack(rBuf, false)]
			log.Error("[session.handleLoop] panic session %s: err=%s\n%s", s.sessionToken(), r, rBuf)
		}

		grNum := atomic.AddInt32(&(s.grNum), -1)
		s.listener.OnClose(s)
		log.Info("%s, [session.handleLoop] goroutine exit now, left gr num %d", s.Stat(), grNum)
		s.gc()
	}()

	flag := true // do not do any read/Write/cron operation while got Write error
	wsConn, wsFlag := s.conn.(*connection.GettyWSConn)
	_, udpFlag := s.conn.(*connection.GettyUDPConn)
	iovec := make([][]byte, 0, maxIovecNum)
	var counter gxtime.CountWatch
LOOP:
	for {
		// A select blocks until one of its cases is ready to run.
		// It choose one at random if multiple are ready. Otherwise it choose default branch if none is ready.
		select {
		case <-s.done:
			// this case assure the (session)handleLoop gr will exit before (session)handlePackage gr.
			<-s.rDone

			if len(s.wQ) == 0 {
				log.Info("%s, [session.handleLoop] got done signal. wQ is nil.", s.Stat())
				break LOOP
			}
			counter.Start()
			if counter.Count() > s.wait.Nanoseconds() {
				log.Info("%s, [session.handleLoop] got done signal ", s.Stat())
				break LOOP
			}

		case outPkg, ok := <-s.wQ:
			if !ok {
				continue
			}
			if !flag {
				log.Warn("[session.handleLoop] drop write out package %#v", outPkg)
				continue
			}

			if udpFlag || wsFlag {
				err := s.WritePkg(outPkg, 0)
				if err != nil {
					log.Error("%s, [session.handleLoop] = error:%+v", s.sessionToken(), errors.WithStack(err))
					s.stop()
					// break LOOP
					flag = false
				}

				continue
			}

			iovec = iovec[:0]
			for idx := 0; idx < maxIovecNum; idx++ {
				pkgBytes, err := s.writer.Write(s, outPkg)
				if err != nil {
					log.Error("%s, [session.handleLoop] = error:%+v", s.sessionToken(), errors.WithStack(err))
					s.stop()
					// break LOOP
					flag = false
					break
				}
				iovec = append(iovec, pkgBytes)

				if idx < maxIovecNum-1 {
					loopFlag := true
					select {
					case outPkg, ok = <-s.wQ:
						if !ok {
							loopFlag = false
						}

					default:
						loopFlag = false
						break
					}
					if !loopFlag {
						break // break for-idx loop
					}
				}
			}
			err := s.WriteBytesArray(iovec[:]...)
			if err != nil {
				log.Error("%s, [session.handleLoop]s.WriteBytesArray(iovec len:%d) = error:%+v",
					s.sessionToken(), len(iovec), errors.WithStack(err))
				s.stop()
				// break LOOP
				flag = false
			}

		case <-wheel.After(s.period):
			if flag {
				if wsFlag {
					err := wsConn.WritePing()
					if err != nil {
						log.Warn("wsConn.writePing() = error:%+v", errors.WithStack(err))
					}
				}
				s.listener.OnCron(s)
			}
		}
	}
}

func (s *session) addTask(pkg interface{}) {
	f := func() {
		s.listener.OnMessage(s, pkg)
		s.conn.IncReadPkgNum()
	}

	// if taskPool := s..GetTaskPool(); taskPool != nil {
	// 	taskPool.AddTask(f)
	// 	return
	// }

	f()
}

func (s *session) handlePackage() {
	var (
		err error
	)

	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			rBuf := make([]byte, size)
			rBuf = rBuf[:runtime.Stack(rBuf, false)]
			log.Error("[session.handlePackage] panic session %s: err=%s\n%s", s.sessionToken(), r, rBuf)
		}

		close(s.rDone)
		grNum := atomic.AddInt32(&(s.grNum), -1)
		log.Info("%s, [session.handlePackage] gr will exit now, left gr num %d", s.sessionToken(), grNum)
		s.stop()
		if err != nil {
			log.Error("%s, [session.handlePackage] error:%+v", s.sessionToken(), errors.WithStack(err))
			if s != nil || s.listener != nil {
				s.listener.OnError(s, err)
			}
		}
	}()

	if _, ok := s.conn.(*connection.GettyTCPConn); ok {
		if s.reader == nil {
			errStr := fmt.Sprintf("session{name:%s, conn:%#v, reader:%#v}", s.name, s.conn, s.reader)
			log.Error(errStr)
			panic(errStr)
		}

		err = s.handleTCPPackage()
	} else if _, ok := s.conn.(*connection.GettyWSConn); ok {
		err = s.handleWSPackage()
	} else if _, ok := s.conn.(*connection.GettyUDPConn); ok {
		err = s.handleUDPPackage()
	} else {
		panic(fmt.Sprintf("unknown type session{%#v}", s))
	}
}

// get package from tcp stream(packet)
func (s *session) handleTCPPackage() error {
	// buf = make([]byte, maxReadBufLen)
	bufp := gxbytes.GetBytes(maxReadBufLen)
	buf := *bufp

	// pktBuf = new(bytes.Buffer)
	pktBuf := gxbytes.GetBytesBuffer()

	defer func() {
		gxbytes.PutBytes(bufp)
		gxbytes.PutBytesBuffer(pktBuf)
	}()

	conn := s.conn.(*connection.GettyTCPConn)
	var err error
	for {
		if s.IsClosed() {
			err = nil
			// do not handle the left stream in pktBuf and exit asap.
			// it is impossible packing a package by the left stream.
			break
		}

		bufLen := 0
		var exit bool
		for {
			// for clause for the network timeout condition check
			// s.conn.SetReadTimeout(time.Now().Add(s.rTimeout))
			bufLen, err = conn.Recv(buf)
			if err != nil {
				if netError, ok := errors.Cause(err).(net.Error); ok && netError.Timeout() {
					break
				}
				if errors.Cause(err) == io.EOF {
					log.Info("%s, [session.conn.read] = error:%+v", s.sessionToken(), errors.WithStack(err))
					err = nil
					exit = true
					break
				}
				log.Error("%s, [session.conn.read] = error:%+v", s.sessionToken(), errors.WithStack(err))
				exit = true
			}
			break
		}
		if exit {
			break
		}
		if 0 == bufLen {
			continue // just continue if session can not read no more stream bytes.
		}
		pktBuf.Write(buf[:bufLen])
		for {
			if pktBuf.Len() <= 0 {
				break
			}
			pkg, pkgLen, err := s.reader.Read(s, pktBuf.Bytes())
			// for case 3/case 4
			if err == nil && s.maxMsgLen > 0 && pkgLen > int(s.maxMsgLen) {
				err = errors.Errorf("pkgLen %d > session max message len %d", pkgLen, s.maxMsgLen)
			}
			// handle case 1
			if err != nil {
				log.Warn("%s, [session.handleTCPPackage] = len{%d}, error:%+v",
					s.sessionToken(), pkgLen, errors.WithStack(err))
				exit = true
				break
			}
			// handle case 2/case 3
			if pkg == nil {
				break
			}
			// handle case 4
			s.conn.UpdateActive()
			s.addTask(pkg)
			pktBuf.Next(pkgLen)
			// continue to handle case 5
		}
		if exit {
			break
		}
	}

	return errors.WithStack(err)
}

// get package from udp packet
func (s *session) handleUDPPackage() error {
	conn := s.conn.(*connection.GettyUDPConn)
	maxBufLen := int(s.maxMsgLen + maxReadBufLen)
	var bufLen int
	if int(s.maxMsgLen<<1) < bufLen {
		maxBufLen = int(s.maxMsgLen << 1)
	}
	bufp := gxbytes.GetBytes(maxBufLen)
	defer gxbytes.PutBytes(bufp)
	buf := *bufp
	var reserr error
	for {
		if s.IsClosed() {
			break
		}

		bufLen, addr, err := conn.Recv(buf)
		log.Debug("conn.read() = bufLen:%d, addr:%#v, err:%+v", bufLen, addr, errors.WithStack(err))
		if netError, ok := errors.Cause(err).(net.Error); ok && netError.Timeout() {
			continue
		}
		if err != nil {
			log.Error("%s, [session.handleUDPPackage] = len:%d, error:%+v",
				s.sessionToken(), bufLen, errors.WithStack(err))
			reserr = errors.Wrapf(err, "conn.read()")
			break
		}

		if bufLen == 0 {
			log.Error("conn.read() = bufLen:%d, addr:%s, err:%+v", bufLen, addr, errors.WithStack(err))
			continue
		}

		if bufLen == len(connectPingPackage) && bytes.Equal(connectPingPackage, buf[:bufLen]) {
			log.Info("got %s connectPingPackage", addr)
			continue
		}

		pkg, pkgLen, err := s.reader.Read(s, buf[:bufLen])
		log.Debug("s.reader.Read() = pkg:%#v, pkgLen:%d, err:%+v", pkg, pkgLen, errors.WithStack(err))
		if err == nil && s.maxMsgLen > 0 && bufLen > int(s.maxMsgLen) {
			err = errors.Errorf("Message Too Long, bufLen %d, session max message len %d", bufLen, s.maxMsgLen)
		}
		if err != nil {
			log.Warn("%s, [session.handleUDPPackage] = len:%d, error:%+v",
				s.sessionToken(), pkgLen, errors.WithStack(err))
			continue
		}
		if pkgLen == 0 {
			log.Error("s.reader.Read() = pkg:%#v, pkgLen:%d, err:%+v", pkg, pkgLen, errors.WithStack(err))
			continue
		}

		s.conn.UpdateActive()
		s.addTask(connection.UDPContext{Pkg: pkg, PeerAddr: addr})
	}

	return errors.WithStack(reserr)
}

// get package from websocket stream
func (s *session) handleWSPackage() error {
	conn := s.conn.(*connection.GettyWSConn)
	for {
		if s.IsClosed() {
			break
		}
		pkg, err := conn.Recv()
		if netError, ok := errors.Cause(err).(net.Error); ok && netError.Timeout() {
			continue
		}
		if err != nil {
			log.Warn("%s, [session.handleWSPackage] = error:%+v",
				s.sessionToken(), errors.WithStack(err))
			return errors.WithStack(err)
		}
		s.conn.UpdateActive()
		if s.reader != nil {
			unmarshalPkg, length, err := s.reader.Read(s, pkg)
			if err == nil && s.maxMsgLen > 0 && length > int(s.maxMsgLen) {
				err = errors.Errorf("Message Too Long, length %d, session max message len %d", length, s.maxMsgLen)
			}
			if err != nil {
				log.Warn("%s, [session.handleWSPackage] = len:%d, error:%+v",
					s.sessionToken(), length, errors.WithStack(err))
				continue
			}

			s.addTask(unmarshalPkg)
		} else {
			s.addTask(pkg)
		}
	}

	return nil
}

func (s *session) stop() {
	select {
	case <-s.done: // s.done is a blocked channel. if it has not been closed, the default branch will be invoked.
		return

	default:
		s.once.Do(func() {
			// let read/Write timeout asap
			now := time.Now()
			if conn := s.Conn(); conn != nil {
				conn.NetConn().SetReadDeadline(now.Add(s.conn.ReadTimeout()))
				conn.NetConn().SetWriteDeadline(now.Add(s.conn.WriteTimeout()))
			}
			close(s.done)
			// c := s.GetAttribute(sessionClientKey)
			// if clt, ok := c.(*client); ok {
			// 	clt.reConnect()
			// }
		})
	}
}

func (s *session) gc() {
	var wQ chan interface{}
	var conn connection.Connection

	s.lock.Lock()
	if s.attrs != nil {
		s.attrs = nil
		if s.wQ != nil {
			wQ = s.wQ
			s.wQ = nil
		}
		conn = s.conn
	}
	s.lock.Unlock()

	go func() {
		if wQ != nil {
			conn.Close((int)(s.wait))
			close(wQ)
		}
	}()
}

// Close will be invoked by NewSessionCallback(if return error is not nil)
// or (session)handleLoop automatically. It's thread safe.
func (s *session) Close() {
	s.stop()
	log.Info("%s closed now. its current gr num is %d",
		s.sessionToken(), atomic.LoadInt32(&(s.grNum)))
}

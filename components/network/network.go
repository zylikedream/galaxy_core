package network

import (
	"time"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/connection"
	"github.com/zylikedream/galaxy/components/network/peer"
)

type NewSessionCallback func(Session) error

// Reader is used to unmarshal a complete pkg from buffer
type Reader interface {
	// Parse tcp/udp/websocket pkg from buffer and if possible return a complete pkg.
	// When receiving a tcp network streaming segment, there are 4 cases as following:
	// case 1: a error found in the streaming segment;
	// case 2: can not unmarshal a pkg header from the streaming segment;
	// case 3: unmarshal a pkg header but can not unmarshal a pkg from the streaming segment;
	// case 4: just unmarshal a pkg from the streaming segment;
	// case 5: unmarshal more than one pkg from the streaming segment;
	//
	// The return value is (nil, 0, error) as case 1.
	// The return value is (nil, 0, nil) as case 2.
	// The return value is (nil, pkgLen, nil) as case 3.
	// The return value is (pkg, pkgLen, nil) as case 4.
	// The handleTcpPackage may invoke func Read many times as case 5.
	Read(Session, []byte) (interface{}, int, error)
}

// Writer is used to marshal pkg and write to session
type Writer interface {
	// if @Session is udpGettySession, the second parameter is UDPContext.
	Write(Session, interface{}) ([]byte, error)
}

// package handler interface
type ReadWriter interface {
	Reader
	Writer
}

// EventListener is used to process pkg that received from remote session
type EventListener interface {
	// invoked when session opened
	// If the return error is not nil, @Session will be closed.
	OnOpen(Session) error

	// invoked when session closed.
	OnClose(Session)

	// invoked when got error.
	OnError(Session, error)

	// invoked periodically, its period can be set by (Session)SetCronPeriod
	OnCron(Session)

	// invoked when getty received a package. Pls attention that do not handle long time
	// logic processing in this func. You'd better set the package's maximum length.
	// If the message's length is greater than it, u should should return err in
	// Reader{Read} and getty will close this connection soon.
	//
	// If ur logic processing in this func will take a long time, u should start a goroutine
	// pool(like working thread pool in cpp) to handle the processing asynchronously. Or u
	// can do the logic processing in other asynchronous way.
	// !!!In short, ur OnMessage callback func should return asap.
	//
	// If this is a udp event listener, the second parameter type is UDPContext.
	OnMessage(Session, interface{})
}

type Session interface {
	Reset()
	Conn() connection.Connection
	Stat() string
	IsClosed() bool
	// get endpoint type
	Peer() peer.Peer

	SetMaxMsgLen(int)
	SetName(string)
	SetEventListener(EventListener)
	SetPkgHandler(ReadWriter)
	SetReader(Reader)
	SetWriter(Writer)
	SetCronPeriod(int)

	SetWQLen(int)
	SetWaitTime(time.Duration)

	GetAttribute(interface{}) interface{}
	SetAttribute(interface{}, interface{})
	RemoveAttribute(interface{})

	// the Writer will invoke this function. Pls attention that if timeout is less than 0, WritePkg will send @pkg asap.
	// for udp session, the first parameter should be UDPContext.
	WritePkg(pkg interface{}, timeout time.Duration) error
	WriteBytes([]byte) error
	WriteBytesArray(...[]byte) error
	Close()
	Run()
}

type Network struct {
	configure *gconfig.Configuration
	peer      peer.Peer
}

func NewNetwork(configFile string) (*Network, error) {
	configure := gconfig.New(configFile)
	peer, err := peer.NewPeer(configure.GetString("network.peer"), configure.WithParent("network"))
	if err != nil {
		return nil, err
	}
	return &Network{
		configure: configure,
		peer:      peer,
	}, nil
}

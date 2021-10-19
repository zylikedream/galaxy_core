package session

import (
	"net"

	"github.com/zylikedream/galaxy/components/network"
	"github.com/zylikedream/galaxy/components/network/connection"
	"github.com/zylikedream/galaxy/components/network/peer"
)

func NewTCPSession(conn net.Conn, peer peer.Peer) network.Session {
	c := connection.NewGettyTCPConn(conn)
	session := newSession(peer, c)
	session.name = defaultTCPSessionName

	return session
}

package session

import (
	"net"

	"github.com/zylikedream/galaxy/components/network"
	"github.com/zylikedream/galaxy/components/network/connection"
	"github.com/zylikedream/galaxy/components/network/peer"
)

func NewUDPSession(conn *net.UDPConn, peer peer.Peer) network.Session {
	c := connection.NewGettyUDPConn(conn)
	session := newSession(peer, c)
	session.name = defaultUDPSessionName

	return session
}

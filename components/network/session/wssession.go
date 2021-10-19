package session

import (
	"github.com/gorilla/websocket"
	"github.com/zylikedream/galaxy/components/network"
	"github.com/zylikedream/galaxy/components/network/connection"
	"github.com/zylikedream/galaxy/components/network/peer"
)

func NewWSSession(conn *websocket.Conn, peer peer.Peer) network.Session {
	c := connection.NewGettyWSConn(conn)
	session := newSession(peer, c)
	session.name = defaultWSSessionName

	return session
}

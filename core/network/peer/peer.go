/*
 * @Author: your name
 * @Date: 2021-11-04 14:34:02
 * @LastEditTime: 2021-11-04 15:13:50
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /components/network/peer/peer.go
 */
package peer

import (
	"context"

	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
	"github.com/zylikedream/galaxy/core/network/message"
	"github.com/zylikedream/galaxy/core/network/session"
)

const (
	PEER_TCP_SERVER    = "peer.tcp_server"
	PEER_TCP_CONNECTOR = "peer.tcp_connector"
)

type Peer interface {
	Start(ctx context.Context, h session.EventHandler) error
	Stop(ctx context.Context)
	GetMessageCodec() message.MessageCodec
}

func NewPeer(t string, c *gconfig.Configuration) (Peer, error) {
	if node, err := gregister.NewNode("peer."+t, c); err != nil {
		return nil, err
	} else {
		return node.(Peer), nil
	}
}

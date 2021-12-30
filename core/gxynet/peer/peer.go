/*
 * @Author: your name
 * @Date: 2021-11-04 14:34:02
 * @LastEditTime: 2021-11-04 15:13:50
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /components/gxynet/peer/peer.go
 */
package peer

import (
	"context"

	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/endpoint"
	"github.com/zylikedream/galaxy/core/gxynet/message"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

const (
	PEER_TCP_SERVER    = "peer.tcp_server"
	PEER_TCP_CONNECTOR = "peer.tcp_connector"
)

type Peer interface {
	Start(ctx context.Context, h endpoint.EventHandler) error
	Stop(ctx context.Context)
	GetMessageCodec() message.MessageCodec
}

func NewPeer(t string, c *gxyconfig.Configuration) (Peer, error) {
	if node, err := gxyregister.NewNode("peer."+t, c); err != nil {
		return nil, err
	} else {
		return node.(Peer), nil
	}
}

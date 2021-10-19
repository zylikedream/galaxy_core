package network

import (
	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/peer"
)

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

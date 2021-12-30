package network

import (
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/network/logger"
	"github.com/zylikedream/galaxy/core/network/peer"
)

type networkConfig struct {
	Peer      string `toml:"peer"`
	LogConfig string `toml:"log_config"`
}

func NewNetwork(configure *gconfig.Configuration) (peer.Peer, error) {
	conf := &networkConfig{}
	if err := configure.UnmarshalKey("network", conf); err != nil {
		return nil, err
	}
	peer, err := peer.NewPeer(conf.Peer, configure)
	if err != nil {
		return nil, err
	}
	if conf.LogConfig == "" {
		logger.SetLogger(gxylog.DefaultLogger())
	} else {
		logger.SetLogger(gxylog.NewLogger("network", gconfig.New(conf.LogConfig)))
	}
	return peer, err
}

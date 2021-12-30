package gxynet

import (
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/gxynet/logger"
	"github.com/zylikedream/galaxy/core/gxynet/peer"
)

type networkConfig struct {
	Peer      string `toml:"peer"`
	LogConfig string `toml:"log_config"`
}

func NewNetwork(configure *gconfig.Configuration) (peer.Peer, error) {
	conf := &networkConfig{}
	if err := configure.UnmarshalKey("gxynet", conf); err != nil {
		return nil, err
	}
	peer, err := peer.NewPeer(conf.Peer, configure)
	if err != nil {
		return nil, err
	}
	if conf.LogConfig == "" {
		logger.SetLogger(gxylog.DefaultLogger())
	} else {
		logger.SetLogger(gxylog.NewLogger("gxynet", gconfig.New(conf.LogConfig)))
	}
	return peer, err
}

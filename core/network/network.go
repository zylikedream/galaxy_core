package network

import (
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/network/logger"
	"github.com/zylikedream/galaxy/core/network/peer"
)

type networkConfig struct {
	Peer      string `toml:"peer"`
	LogConfig string `toml:"log_config"`
}

func NewNetwork(configFile string) (peer.Peer, error) {
	configure := gconfig.New(configFile)
	conf := &networkConfig{}
	if err := configure.UnmarshalKey("network", conf); err != nil {
		return nil, err
	}
	peer, err := peer.NewPeer(conf.Peer, configure)
	if err != nil {
		return nil, err
	}
	if conf.LogConfig == "" {
		logger.SetLogger(glog.DefaultLogger())
	} else {
		logger.SetLogger(glog.NewLogger("network", conf.LogConfig))
	}
	return peer, err
}

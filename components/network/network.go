package network

import (
	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/glog"
	"github.com/zylikedream/galaxy/components/network/logger"
	"github.com/zylikedream/galaxy/components/network/peer"
)

type networkConfig struct {
	PeerType  string `toml:"peer_type"`
	LogConfig string `toml:"log_config"`
}

func NewNetwork(configFile string) (peer.Peer, error) {
	configure := gconfig.New(configFile)
	conf := &networkConfig{}
	if err := configure.UnmarshalKey("network", conf); err != nil {
		return nil, err
	}
	peer, err := peer.NewPeer(conf.PeerType, configure.WithParent("network"))
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

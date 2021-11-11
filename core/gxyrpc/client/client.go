package client

import (
	"fmt"

	"github.com/smallnest/rpcx/client"
	"github.com/zylikedream/galaxy/core/gconfig"
	discovery "github.com/zylikedream/galaxy/core/gxyrpc/client/Discovery"
)

type gxyrpcClient struct {
	conf *clientConfig
	pool *client.OneClientPool
}

type clientConfig struct {
	PoolSize          int    `toml:"pool_size"`
	FailMode          string `toml:"fail_mode"`
	DefaultSelectMode string `toml:"default_select_mode"`
	Discovery         string `toml:"discovery"`
}

func parseFailMode(mode string) client.FailMode {
	switch mode {
	case "fail_over":
		return client.Failover
	case "fail_fast":
		return client.Failfast
	case "fail_try":
		return client.Failtry
	case "fail_backup":
		return client.Failbackup
	}
	return -1
}

func parseSelectMode(mode string) client.SelectMode {
	switch mode {
	case "random_select":
		return client.RandomSelect
	case "roundrobin":
		return client.RoundRobin
	case "weighted_roundrobin":
		return client.WeightedRoundRobin
	case "weighted_icmp":
		return client.WeightedICMP
	case "consthash":
		return client.ConsistentHash
	case "closest":
		return client.Closest
	case "custom":
		return client.SelectByUser
	}
	return -1
}

func NewGrpcClient(configFile string) (*gxyrpcClient, error) {
	conf := &clientConfig{}
	configure := gconfig.New(configFile)
	if err := configure.UnmarshalKey("gxyrpc", conf); err != nil {
		return nil, err
	}
	gxyrpc := &gxyrpcClient{
		conf: conf,
	}
	d, err := discovery.NewDisvoery(conf.Discovery, configure)
	if err != nil {
		return nil, err
	}
	failMode := parseFailMode(conf.FailMode)
	if failMode < 0 {
		return nil, fmt.Errorf("unkown fail mode")
	}
	selectMode := parseSelectMode(conf.DefaultSelectMode)
	if selectMode < 0 {
		return nil, fmt.Errorf("unkonw select mode")
	}
	gxyrpc.pool = client.NewOneClientPool(conf.PoolSize, failMode, selectMode, d.GetDiscovery(), client.DefaultOption)
	return gxyrpc, nil
}

package gxyrpc

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/smallnest/rpcx/client"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyrpc/client/discovery"
)

type RpcClient struct {
	conf           *clientConfig
	defaultClient  *client.OneClient
	serviceClients map[string]client.XClient
}

type clientConfig struct {
	Default  serviceConfig   `toml:"default"`
	Services []serviceConfig `toml:"services"`
}

type serviceConfig struct {
	Service    string `toml:"service"`
	FailMode   string `toml:"fail_mode"`
	SelectMode string `toml:"select_mode"`
	Discovery  string `toml:"discovery"`
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

func NewGrpcClient(configFile string) (*RpcClient, error) {
	conf := &clientConfig{}
	configure := gxyconfig.New(configFile)
	if err := configure.UnmarshalKey("gxyrpc_client", conf); err != nil {
		return nil, err
	}
	gxyrpc := &RpcClient{
		conf:           conf,
		serviceClients: make(map[string]client.XClient),
	}
	defaultClient, err := newDefaultServiceClient(&conf.Default, configure)
	if err != nil {
		return nil, err
	}
	gxyrpc.defaultClient = defaultClient
	for _, serviceConf := range conf.Services {
		cli, err := newServiceClient(&serviceConf, configure)
		if err != nil {
			return nil, err
		}
		gxyrpc.serviceClients[serviceConf.Service] = cli
	}
	return gxyrpc, nil
}

func newDefaultServiceClient(sc *serviceConfig, c *gxyconfig.Configuration) (*client.OneClient, error) {
	d, err := discovery.NewDisvoery(sc.Discovery, c)
	if err != nil {
		return nil, err
	}
	failMode := parseFailMode(sc.FailMode)
	if failMode < 0 {
		return nil, fmt.Errorf("unkown fail mode:%s", sc.FailMode)
	}
	selectMode := parseSelectMode(sc.SelectMode)
	if selectMode < 0 {
		return nil, fmt.Errorf("unkonw select mode:%s", sc.SelectMode)
	}
	return client.NewOneClient(failMode, selectMode, d.GetDiscovery(), client.DefaultOption), nil

}

func newServiceClient(sc *serviceConfig, c *gxyconfig.Configuration) (client.XClient, error) {
	d, err := discovery.NewDisvoery(sc.Discovery, c)
	if err != nil {
		return nil, err
	}
	failMode := parseFailMode(sc.FailMode)
	if failMode < 0 {
		return nil, fmt.Errorf("unkown fail mode:%s", sc.FailMode)
	}
	selectMode := parseSelectMode(sc.SelectMode)
	if selectMode < 0 {
		return nil, fmt.Errorf("unkonw select mode:%s", sc.SelectMode)
	}
	return client.NewXClient(sc.Service, failMode, selectMode, d.GetDiscovery(), client.DefaultOption), nil

}

func (g *RpcClient) Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (*client.Call, error) {
	cli, ok := g.serviceClients[servicePath]
	if ok {
		return cli.Go(ctx, serviceMethod, args, reply, done)
	}
	return g.defaultClient.Go(ctx, servicePath, serviceMethod, args, reply, done)

}

func (g *RpcClient) Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) error {
	cli, ok := g.serviceClients[servicePath]
	if ok {
		return cli.Call(ctx, serviceMethod, args, reply)
	}
	return g.defaultClient.Call(ctx, servicePath, serviceMethod, args, reply)
}

func (g *RpcClient) Fork(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) error {
	cli, ok := g.serviceClients[servicePath]
	if ok {
		return cli.Fork(ctx, serviceMethod, args, reply)
	}
	return g.defaultClient.Fork(ctx, servicePath, serviceMethod, args, reply)
}

func (g *RpcClient) Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) error {
	cli, ok := g.serviceClients[servicePath]
	if ok {
		return cli.Broadcast(ctx, serviceMethod, args, reply)
	}
	return g.defaultClient.Broadcast(ctx, servicePath, serviceMethod, args, reply)
}

func (g *RpcClient) SendFile(ctx context.Context, fileName string, rateInBytesPerSecond int64, meta map[string]string) error {
	return g.defaultClient.SendFile(ctx, fileName, rateInBytesPerSecond, meta) // file的service是固定的 直接使用defaultclient即可
}

func (g *RpcClient) DownloadFile(ctx context.Context, requestFileName string, saveTo io.Writer, meta map[string]string) error {
	return g.defaultClient.DownloadFile(ctx, requestFileName, saveTo, meta) // file的service是固定的 直接使用defaultclient即可
}

func (g *RpcClient) Stream(ctx context.Context, meta map[string]string) (net.Conn, error) {
	return g.defaultClient.Stream(ctx, meta)
}

func (g *RpcClient) Close() error {
	var errs []error
	if err := g.defaultClient.Close(); err != nil {
		errs = append(errs, err)
	}
	for _, serviceCli := range g.serviceClients {
		if err := serviceCli.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		// todo 整合错误信息
		return errs[0]
	}
	return nil
}

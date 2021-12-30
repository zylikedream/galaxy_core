package registery

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type peerRegistery struct {
	conf *peerConfig
}

type peerConfig struct {
}

func newPeerRegistery(_ *gconfig.Configuration) (*peerRegistery, error) {
	regist := &peerRegistery{
		conf: &peerConfig{},
	}
	return regist, nil
}

func (r *peerRegistery) Type() string {
	return REGISTERY_TYPE_PEER
}

func (r *peerRegistery) Start() error {
	return nil
}

func (r *peerRegistery) GetPlugin() server.Plugin {
	return nil
}

func (t *peerRegistery) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newPeerRegistery(c)
}

func init() {
	gxyregister.Register((*peerRegistery)(nil))
}

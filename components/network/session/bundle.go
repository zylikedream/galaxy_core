package session

import (
	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/processor"
)

type SessionBundle struct {
	Proc    *processor.Processor
	Handler EventHandler
}

func (p *SessionBundle) BindProc(c *gconfig.Configuration) error {
	proc, err := processor.NewProcessor(c)
	if err != nil {
		return err
	}
	p.Proc = proc
	return nil
}

func (p *SessionBundle) BindHandler(handler EventHandler) {
	p.Handler = handler
}

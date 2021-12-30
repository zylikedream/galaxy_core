package session

import (
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxynet/processor"
)

type SessionBundle struct {
	*processor.Processor
	Handler EventHandler
}

func (p *SessionBundle) BindProc(c *gconfig.Configuration) error {
	proc, err := processor.NewProcessor(c)
	if err != nil {
		return err
	}
	p.Processor = proc
	return nil
}

func (p *SessionBundle) BindHandler(handler EventHandler) {
	p.Handler = handler
}

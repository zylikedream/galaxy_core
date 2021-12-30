package session

import (
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/processor"
)

type SessionBundle struct {
	*processor.Processor
	Handler EventHandler
}

func (p *SessionBundle) BindProc(c *gxyconfig.Configuration) error {
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

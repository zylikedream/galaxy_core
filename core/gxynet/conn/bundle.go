package conn

import (
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/processor"
)

type ConnBundle struct {
	*processor.Processor
	Handler EventHandler
}

func (p *ConnBundle) BindProc(c *gxyconfig.Configuration) error {
	proc, err := processor.NewProcessor(c)
	if err != nil {
		return err
	}
	p.Processor = proc
	return nil
}

func (p *ConnBundle) BindHandler(handler EventHandler) {
	p.Handler = handler
}

package endpoint

import (
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/processor"
)

type CoreBundle struct {
	processor.Processor
	Handler EventHandler
}

func (p *CoreBundle) BindProc(c *gxyconfig.Configuration, procType string) error {
	proc, err := processor.NewProcessor(procType, c)
	if err != nil {
		return err
	}
	p.Processor = proc
	return nil
}

func (p *CoreBundle) BindHandler(handler EventHandler) {
	p.Handler = handler
}

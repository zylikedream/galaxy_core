package processor

import "github.com/zylikedream/galaxy/components/gconfig"

type ProcessorBundle struct {
	Proc    *Processor
	Handler MsgHandler
}

func (p *ProcessorBundle) BindProc(c *gconfig.Configuration) error {
	proc, err := NewProcessor(c)
	if err != nil {
		return err
	}
	p.Proc = proc
	return nil
}

func (p *ProcessorBundle) BindHandler(handler MsgHandler) {
	p.Handler = handler
}

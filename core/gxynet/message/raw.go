package message

import (
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type RawMessage struct {
}

func newRawMessage(_ *gxyconfig.Configuration) (*RawMessage, error) {
	return &RawMessage{}, nil
}

func (j *RawMessage) Decode(msg interface{}, data []byte) error {
	return nil
}

func (j *RawMessage) Encode(msg interface{}) ([]byte, error) {
	return msg.([]byte), nil
}

func (j *RawMessage) Type() string {
	return MESSAGE_JSON
}

func (j *RawMessage) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newRawMessage(c)
}

func init() {
	gxyregister.Register((*RawMessage)(nil))
}

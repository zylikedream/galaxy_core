package message

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type JsonMessage struct {
}

func newJsonMessage(_ *gxyconfig.Configuration) (*JsonMessage, error) {
	return &JsonMessage{}, nil
}

func (j *JsonMessage) Decode(msg interface{}, data []byte) error {
	if err := json.Unmarshal(data, msg); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (j *JsonMessage) Encode(msg interface{}) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

func (j *JsonMessage) Type() string {
	return MESSAGE_JSON
}

func (j *JsonMessage) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newJsonMessage(c)
}

func init() {
	gxyregister.Register((*JsonMessage)(nil))
}

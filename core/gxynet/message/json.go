package message

import (
	"encoding/json"

	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type JsonMessage struct {
}

func newJsonMessage(_ *gxyconfig.Configuration) (*JsonMessage, error) {
	return &JsonMessage{}, nil
}

func (j *JsonMessage) Decode(msg interface{}, data []byte) error {
	return json.Unmarshal(data, msg)
}

func (j *JsonMessage) Encode(msg interface{}) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
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

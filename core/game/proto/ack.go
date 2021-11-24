package proto

import (
	"hash/crc32"

	"github.com/zylikedream/galaxy/core/network/message"
)

type Ack struct {
	Code   int    `json:"code"`
	MsgID  int    `json:"msg_id"`
	Reason string `json:"reason,omitempty"`
	Data   []byte `json:"data,omitempty"`
}

func init() {
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("Ack"))), (*Ack)(nil))
}

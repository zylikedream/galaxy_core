package proto

import (
	"hash/crc32"

	"github.com/zylikedream/galaxy/components/network/message"
)

type EchoReq struct {
	Msg string `json:"msg"`
}

type EchoAck struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func init() {
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("EchoReq"))), (*EchoReq)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("EchoAck"))), (*EchoAck)(nil))
}

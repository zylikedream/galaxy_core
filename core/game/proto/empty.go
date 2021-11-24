package proto

import (
	"hash/crc32"

	"github.com/zylikedream/galaxy/core/network/message"
)

type RspEmpty struct {
}

func init() {
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("RspEmpty"))), (*RspEmpty)(nil))
}

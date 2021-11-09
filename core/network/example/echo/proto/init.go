package proto

import (
	"hash/crc32"

	"github.com/zylikedream/galaxy/core/network/message"
)

func init() {
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("EchoReq"))), (*EchoReq)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("EchoAck"))), (*EchoAck)(nil))
}

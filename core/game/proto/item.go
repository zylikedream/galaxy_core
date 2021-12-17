package proto

import (
	"hash/crc32"

	"github.com/zylikedream/galaxy/core/network/message"
)

// 账号登录
type NtfItemUpdate struct {
	Items []PItemInfo `json:"items"`
}

// 物品信息
type PItemInfo struct {
	PropID int    `json:"prop_id"`
	Num    uint64 `json:"num"`
}

func init() {
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("NtfItemUpdate"))), (*NtfItemUpdate)(nil))
}

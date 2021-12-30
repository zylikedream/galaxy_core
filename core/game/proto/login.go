package proto

import (
	"hash/crc32"

	"github.com/zylikedream/galaxy/core/gxynet/message"
)

type ReqHandshake struct {
	LoginKey string `json:"login_key"`
}

type RspHandshake struct {
	Timestamp uint64 `json:"timestamp"`
}

// 账号登录
type ReqAccountLogin struct {
	Account    string      `json:"account"`
	ClientInfo PClientInfo `json:"client_info"`
}

// 客户端信息
type PClientInfo struct {
	SdkType    int    `json:"sdk_type"`
	SysVersion string `json:"sys_version"`
	DevID      string `json:"dev_id"`
}

type RspAccountLogin struct {
	Create bool `json:"create"`
}

func init() {
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("ReqHandshake"))), (*ReqHandshake)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("RspHandshake"))), (*RspHandshake)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("ReqAccountLogin"))), (*ReqAccountLogin)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("RspAccountLogin"))), (*RspAccountLogin)(nil))
}

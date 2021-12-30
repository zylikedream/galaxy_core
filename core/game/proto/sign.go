package proto

import (
	"hash/crc32"

	"github.com/zylikedream/galaxy/core/gxynet/message"
)

type ReqSignInfo struct {
}

type RspSignInfo struct {
	SignTime int64 `json:"sign_time"`
	SignDay  int   `json:"sign_day"`
	Day      int   `json:"day"`
}

type ReqSignSign struct {
}

type RspSignSign struct {
}

type ReqSignRepair struct {
	SignDay int
}

type RspSignRepair struct {
}

func init() {
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("reqSignInfo"))), (*ReqSignInfo)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("RspSignInfo"))), (*RspSignInfo)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("reqSignSign"))), (*ReqSignSign)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("rspSignSign"))), (*RspSignSign)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("reqSignRepair"))), (*ReqSignRepair)(nil))
	message.RegisterMessageMeta(int(crc32.ChecksumIEEE([]byte("rspSignRepair"))), (*RspSignRepair)(nil))
}

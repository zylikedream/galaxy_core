package processor

import "github.com/zylikedream/galaxy/components/network/message"

type MsgHandler func(msg *message.Message) error

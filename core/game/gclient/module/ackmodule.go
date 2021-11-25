package module

import (
	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/network/message"
	"go.uber.org/zap"
)

type AckModule struct {
	BaseModule
}

func (l *AckModule) Ack(ctx gcontext.Context, ack *proto.Ack) error {
	meta := message.MessageMetaByID(ack.MsgID)
	if ack.Code != ACK_CODE_OK {
		glog.Error("ack failed", zap.String("msg", meta.TypeName()), zap.String("reason", ack.Reason))
		return nil
	}
	sess := GetSessionFromCtx(ctx)
	msg := meta.NewInstance()
	if err := sess.GetMessageCodec().Decode(msg, ack.Data); err != nil {
		return err
	}
	if err := HandleMessage(ctx, msg); err != nil {
		return err
	}
	glog.Debug("ack success", zap.String("msg", meta.TypeName()))
	return nil
}

func init() {
	Register(&AckModule{})
}

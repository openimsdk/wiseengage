package customerservice

import (
	"context"
	"time"

	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (o *customerService) CallbackAfterSendSingleMsgCommand(ctx context.Context, req *customerservice.CallbackAfterSendSingleMsgCommandReq) (*customerservice.CallbackAfterSendSingleMsgCommandResp, error) {
	lastMsg := model.LastMessage{
		MsgID:      req.MsgID,
		SendTime:   time.UnixMilli(req.SendTime),
		UserID:     req.UserID,
		UpdateTime: time.Now(),
	}
	if err := o.db.UpdateConversationLastMsg(ctx, req.ConversationID, &lastMsg); err != nil {
		return nil, err
	}
	// todo: send to ai
	return &customerservice.CallbackAfterSendSingleMsgCommandResp{}, nil
}

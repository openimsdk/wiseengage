package customerapi

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/wiseengage/v1/pkg/callbackstruct"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (x *API) CallbackAfterSendSingleMsgCommand(c *gin.Context) {
	var cbReq callbackstruct.CallbackAfterSendSingleMsgReq
	if err := c.BindJSON(&cbReq); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if cbReq.SessionType != constant.ReadGroupChatType {
		apiresp.GinSuccess(c, nil)
		return
	}
	rpcReq := &customerservice.CallbackAfterSendSingleMsgCommandReq{
		ConversationID: "sg_" + cbReq.RecvID,
		UserID:         cbReq.SendID,
		MsgID:          cbReq.ServerMsgID,
		SendTime:       cbReq.SendTime,
		Content:        cbReq.Content,
		ContentType:    cbReq.ContentType,
	}
	if _, err := x.client.CallbackAfterSendSingleMsgCommand(c, rpcReq); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

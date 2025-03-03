package customerservice

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (x *API) OpenIMCallback(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	command := strings.TrimPrefix(c.Param(constant.CallbackCommand), "/")
	if command == "callbackAfterSendGroupMsgCommand" {
		x.callbackAfterSendGroupMsg(c, body)
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	log.ZWarn(c, "OpenIMCallback unknown command", nil, "command", command, "body", string(body))
	c.JSON(http.StatusOK, gin.H{})
}

func (x *API) callbackAfterSendGroupMsg(ctx context.Context, body []byte) {
	type msgKeyField struct {
		ServerMsgID string `json:"serverMsgID"`
		SendID      string `json:"sendID"`
		GroupID     string `json:"groupID"`
		SendTime    int64  `json:"sendTime"`
	}
	var msgKey msgKeyField
	if err := json.Unmarshal(body, &msgKey); err != nil {
		log.ZError(ctx, "callbackAfterSendGroupMsg unmarshal failed", err, "body", string(body))
		return
	}
	if msgKey.ServerMsgID == "" || msgKey.SendID == "" || msgKey.GroupID == "" {
		log.ZDebug(ctx, "callbackAfterSendGroupMsg invalid msgKey", "body", string(body))
		return
	}
	req := &customerservice.UpdateSendMsgTimeReq{
		ConversationID: "sg_" + msgKey.GroupID,
		UserID:         msgKey.SendID,
		SendTime:       msgKey.SendTime,
		MsgID:          msgKey.ServerMsgID,
	}
	if _, err := x.client.UpdateSendMsgTime(ctx, req); err != nil {
		log.ZError(ctx, "callbackAfterSendGroupMsg UpdateSendMsgTime failed", err, "req", req)
	}
}

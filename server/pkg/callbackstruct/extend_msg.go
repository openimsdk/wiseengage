package callbackstruct

import (
	"github.com/openimsdk/protocol/extendmsg"
)

type CallbackBeforeSetMessageReactionExtReq struct {
	OperationID           string `json:"operationID"`
	CallbackCommand       `json:"callbackCommand"`
	ConversationID        string                         `json:"conversationID"`
	OpUserID              string                         `json:"opUserID"`
	SessionType           int32                          `json:"sessionType"`
	ReactionExtensionList map[string]*extendmsg.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string                         `json:"clientMsgID"`
	IsReact               bool                           `json:"isReact"`
	IsExternalExtensions  bool                           `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                          `json:"msgFirstModifyTime"`
}
type CallbackBeforeSetMessageReactionExtResp struct {
	CommonCallbackResp
	ResultReactionExtensionList []*extendmsg.KeyValueResp `json:"resultReactionExtensionList"`
	MsgFirstModifyTime          int64                     `json:"msgFirstModifyTime"`
}
type CallbackDeleteMessageReactionExtReq struct {
	CallbackCommand       `json:"callbackCommand"`
	OperationID           string                `json:"operationID"`
	ConversationID        string                `json:"conversationID"`
	OpUserID              string                `json:"opUserID"`
	SessionType           int32                 `json:"sessionType"`
	ReactionExtensionList []*extendmsg.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string                `json:"clientMsgID"`
	IsExternalExtensions  bool                  `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                 `json:"msgFirstModifyTime"`
}
type CallbackDeleteMessageReactionExtResp struct {
	CommonCallbackResp
	ResultReactionExtensionList []*extendmsg.KeyValueResp `json:"resultReactionExtensionList"`
	MsgFirstModifyTime          int64                     `json:"msgFirstModifyTime"`
}

type CallbackGetMessageListReactionExtReq struct {
	OperationID     string `json:"operationID"`
	CallbackCommand `json:"callbackCommand"`
	ConversationID  string   `json:"conversationID"`
	OpUserID        string   `json:"opUserID"`
	SessionType     int32    `json:"sessionType"`
	TypeKeyList     []string `json:"typeKeyList"`
	//MessageKeyList  []*msg.GetMessageListReactionExtensionsReq_MessageReactionKey `json:"messageKeyList"`
}

type CallbackGetMessageListReactionExtResp struct {
	CommonCallbackResp
	MessageResultList []*extendmsg.SingleMessageExtensionResult `json:"messageResultList"`
}

type CallbackAddMessageReactionExtReq struct {
	OperationID           string `json:"operationID"`
	CallbackCommand       `json:"callbackCommand"`
	ConversationID        string                         `json:"conversationID"`
	OpUserID              string                         `json:"opUserID"`
	SessionType           int32                          `json:"sessionType"`
	ReactionExtensionList map[string]*extendmsg.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string                         `json:"clientMsgID"`
	IsReact               bool                           `json:"isReact"`
	IsExternalExtensions  bool                           `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                          `json:"msgFirstModifyTime"`
}

type CallbackAddMessageReactionExtResp struct {
	CommonCallbackResp
	ResultReactionExtensionList []*extendmsg.KeyValueResp `json:"resultReactionExtensionList"`
	IsReact                     bool                      `json:"isReact"`
	MsgFirstModifyTime          int64                     `json:"msgFirstModifyTime"`
}

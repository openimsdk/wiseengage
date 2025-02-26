package imapi

import "github.com/openimsdk/protocol/sdkws"

type SendBusinessNotificationReq struct {
	Key        string `json:"key"`
	Data       string `json:"data"`
	SendUserID string `json:"sendUserID"`
	RecvUserID string `json:"recvUserID"`
}

type SendBusinessNotificationResp struct {
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
}

type SendMsgReq struct {
	RecvID           string                 `json:"recvID"`
	SendID           string                 `json:"sendID"`
	GroupID          string                 `json:"groupID"`
	SenderNickname   string                 `json:"senderNickname"`
	SenderFaceURL    string                 `json:"senderFaceURL"`
	SenderPlatformID int32                  `json:"senderPlatformID"`
	Content          map[string]any         `json:"content"`
	ContentType      int32                  `json:"contentType"`
	SessionType      int32                  `json:"sessionType"`
	IsOnlineOnly     bool                   `json:"isOnlineOnly"`
	NotOfflinePush   bool                   `json:"notOfflinePush"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}

type SendMsgResp struct {
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
}

package customerservice

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"time"

	constant2 "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/idutil"
	"github.com/openimsdk/tools/utils/timeutil"
	"github.com/openimsdk/wiseengage/v1/pkg/common/config"
	"github.com/openimsdk/wiseengage/v1/pkg/common/constant"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/controller"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	pb "github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (o *customerService) initConversation(ctx context.Context, userID string) (*model.Conversation, error) {
	if conversation, err := o.db.TakeConversationByUserID(ctx, userID); err == nil {
		return conversation, nil
	} else if !controller.IsNotFound(err) {
		return nil, err
	}
	joinGroupResp, err := o.groupCli.GetFullJoinGroupIDs(ctx, &group.GetFullJoinGroupIDsReq{UserID: userID})
	if err != nil {
		return nil, err
	}
	var groupID string
	if len(joinGroupResp.GroupIDs) > 0 {
		groupID = joinGroupResp.GroupIDs[0]
	} else {
		req := &group.CreateGroupReq{OwnerUserID: o.robotUserID, MemberUserIDs: []string{userID}, GroupInfo: &sdkws.GroupInfo{GroupID: "wiseengage"}}
		createGroupResp, err := o.groupCli.CreateGroup(ctx, req)
		if err != nil {
			return nil, err
		}
		groupID = createGroupResp.GroupInfo.GroupID
	}
	now := time.Now()
	conversation := &model.Conversation{
		UserID:         userID,
		ConversationID: "sg_" + groupID,
		CreateTime:     now,
		Status:         constant.ConversationStatusClosed,
	}
	if err := o.db.CreateConversation(ctx, conversation); err != nil {
		if controller.IsDuplicateKeyError(err) {
			return o.db.TakeConversationByUserID(ctx, userID)
		}
		return nil, err
	}
	return conversation, nil
}

func (o *customerService) StartConsultation(ctx context.Context, req *pb.StartConsultationReq) (*pb.StartConsultationResp, error) {
	// TODO: check user
	conversation, err := o.initConversation(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if conversation.Status == constant.ConversationStatusClosed {
		updated, err := o.db.UpdateConversationStatusOpen(ctx, req.UserID,
			conversation.ConversationID, conversation.Version, constant.ConversationRoleRobot)
		if err != nil {
			return nil, err
		}
		if updated {
			if err := o.sendMsg(ctx, conversation.ConversationID, o.config.Msg.Start); err != nil {
				return nil, err
			}
		}
	}
	return &pb.StartConsultationResp{ConversationID: conversation.ConversationID}, nil
}

func (o *customerService) sendMsg(ctx context.Context, conversationID string, msgs []config.Message) error {
	if len(msgs) == 0 {
		return nil
	}
	groupID := strings.TrimPrefix(conversationID, "sg_")
	now := timeutil.GetCurrentTimestampByMill()
	req := msg.SendMsgReq{
		MsgData: &sdkws.MsgData{
			SendID:           o.robotUserID,
			GroupID:          groupID,
			ClientMsgID:      idutil.GetMsgIDByMD5(strconv.Itoa(int(rand.Uint32()))),
			SenderPlatformID: constant2.AdminPlatformID,
			SenderNickname:   "robot",
			SenderFaceURL:    "",
			SessionType:      constant2.ReadGroupChatType,
			MsgFrom:          constant2.SysMsgType,
			ContentType:      0,
			Content:          nil,
			CreateTime:       now,
			SendTime:         now,
			Options: map[string]bool{
				constant2.IsHistory:            false,
				constant2.IsPersistent:         false,
				constant2.IsSenderSync:         false,
				constant2.IsConversationUpdate: false,
			},
			OfflinePushInfo: nil,
			Ex:              "",
		},
	}
	for _, m := range msgs {
		req.MsgData.ContentType = m.ContentType
		req.MsgData.Content = []byte(m.Content)
		if _, err := o.msgCli.SendMsg(ctx, &req); err != nil {
			return err
		}
	}
	return nil
}

func (o *customerService) UpdateSendMsgTime(ctx context.Context, req *pb.UpdateSendMsgTimeReq) (*pb.UpdateSendMsgTimeResp, error) {
	lastMsg := &model.LastMessage{
		Seq:        req.SendMsgSeq,
		SendTime:   time.UnixMilli(req.SendMsgTime),
		UserID:     req.SendUserID,
		UpdateTime: time.Now(),
	}
	if err := o.db.UpdateConversationLastMsg(ctx, req.UserID, req.ConversationID, lastMsg); err != nil {
		return nil, err
	}
	return &pb.UpdateSendMsgTimeResp{}, nil
}

func (o *customerService) UpdateConversationClosed(ctx context.Context, req *pb.UpdateConversationClosedReq) (*pb.UpdateConversationClosedResp, error) {
	updated, err := o.db.UpdateConversationStatusClosed(ctx, req.UserID, req.ConversationID, int(req.Version), req.Cause)
	if err != nil {
		return nil, err
	}
	if updated {
		if req.Timeout {
			err = o.sendMsg(ctx, req.ConversationID, o.config.Msg.Timeout)
		} else {
			err = o.sendMsg(ctx, req.ConversationID, o.config.Msg.Closed)
		}
		if err != nil {
			log.ZError(ctx, "send closed msg", err, "req", req)
		}
	}
	return &pb.UpdateConversationClosedResp{}, nil
}

func (o *customerService) ChangeConversationRole(ctx context.Context, req *pb.ChangeConversationRoleReq) (*pb.ChangeConversationRoleResp, error) {
	updated, err := o.db.UpdateConversationRole(ctx, req.UserID, req.ConversationID, -1, req.Role)
	if err != nil {
		return nil, err
	}
	if updated {
		// todo: change role
	}
	return &pb.ChangeConversationRoleResp{}, nil
}

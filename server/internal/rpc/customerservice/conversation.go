package customerservice

import (
	"context"
	"time"

	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
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
		updated, err := o.db.UpdateConversationStatus(ctx, req.UserID, conversation.ConversationID, constant.ConversationStatusOpen, constant.ConversationRoleRobot, conversation.Version)
		if err != nil {
			return nil, err
		}
		if updated {
			// todo send welcome message
		}
	}
	return &pb.StartConsultationResp{ConversationID: conversation.ConversationID}, nil
}

func (o *customerService) UpdateSendMsgTime(ctx context.Context, req *pb.UpdateSendMsgTimeReq) (*pb.UpdateSendMsgTimeResp, error) {
	lastMsg := &model.LastMessage{
		Seq:      req.SendMsgSeq,
		SendTime: time.UnixMilli(req.SendMsgTime),
		UserID:   req.SendUserID,
	}
	if err := o.db.UpdateConversationLastMsg(ctx, req.UserID, req.ConversationID, lastMsg); err != nil {
		return nil, err
	}
	return &pb.UpdateSendMsgTimeResp{}, nil
}

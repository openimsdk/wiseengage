package customerservice

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand/v2"
	"strings"
	"time"

	constant2 "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/wiseengage/v1/pkg/common/constant"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/controller"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"github.com/openimsdk/wiseengage/v1/pkg/imapi"
	pb "github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (o *customerService) initConversation(ctx context.Context, userID string, agentUserID string) (*model.Conversation, error) {
	if conversation, err := o.db.TakeConversationByUserID(ctx, userID); err == nil {
		return conversation, nil
	} else if !controller.IsNotFound(err) {
		return nil, err
	}
	joinGroupIDs, err := o.imApi.GetFullJoinGroupIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	var groupID string
	if len(joinGroupIDs) > 0 {
		groupID = joinGroupIDs[0]
	} else {
		req := &group.CreateGroupReq{OwnerUserID: agentUserID, MemberUserIDs: []string{userID}, GroupInfo: &sdkws.GroupInfo{GroupID: "wiseengage"}}
		groupInfo, err := o.imApi.CreateGroup(ctx, req)
		if err != nil {
			return nil, err
		}
		groupID = groupInfo.GroupID
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
	agents, err := o.db.FindAgentType(ctx, req.AgentType, []string{constant.AgentStatusEnable})
	if err != nil {
		return nil, err
	}
	if len(agents) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("no available agent")
	}
	conversation, err := o.initConversation(ctx, req.UserID, agents[rand.Uint32()%uint32(len(agents))].UserID)
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
			o.agentSendConfigMessage(ctx, conversation.ConversationID, conversation.AgentUserID, func(agent *model.Agent) *model.AgentMessage {
				return agent.StartMsg
			})
		}
	}
	return &pb.StartConsultationResp{ConversationID: conversation.ConversationID}, nil
}

func (o *customerService) agentSendConfigMessage(ctx context.Context, conversationID string, agentUserID string, msgFn func(agent *model.Agent) *model.AgentMessage) {
	if agentUserID == "" {
		conversation, err := o.db.TakeConversationByUserID(ctx, conversationID)
		if err != nil {
			log.ZWarn(ctx, "send msg take conversation", err, "conversationID", conversationID)
			return
		}
		agentUserID = conversation.AgentUserID
	}
	agent, err := o.db.TakeAgent(ctx, agentUserID)
	if err != nil {
		log.ZWarn(ctx, "sendMsg take agent", err, "agentUserID", agentUserID)
		return
	}
	val := msgFn(agent)
	if val == nil {
		return
	}
	decoder := json.NewDecoder(bytes.NewReader([]byte(val.Content)))
	decoder.UseNumber()
	var message map[string]any
	if err := decoder.Decode(&message); err != nil {
		log.ZError(ctx, "agentSendConfigMessage decode", err, "conversationID", conversationID, "message", val)
		return
	}
	if err := o.agentSendMsg(ctx, conversationID, agent, val.ContentType, message); err != nil {
		log.ZWarn(ctx, "sendMsg take agent", err, "agentUserID", agentUserID)
	}
}

func (o *customerService) agentSendMsg(ctx context.Context, conversationID string, agent *model.Agent, contentType int32, content map[string]any) error {
	req := &imapi.SendMsgReq{
		SendID:           agent.UserID,
		GroupID:          strings.TrimPrefix(conversationID, "sg_"),
		SenderNickname:   agent.Nickname,
		SenderFaceURL:    agent.FaceURL,
		SenderPlatformID: constant2.AdminPlatformID,
		SessionType:      constant2.ReadGroupChatType,
		ContentType:      contentType,
		Content:          content,
	}
	_, err := o.imApi.SendMsg(ctx, req)
	return err
}

func (o *customerService) UpdateSendMsgTime(ctx context.Context, req *pb.UpdateSendMsgTimeReq) (*pb.UpdateSendMsgTimeResp, error) {
	lastMsg := &model.LastMessage{
		MsgID:      req.MsgID,
		SendTime:   time.UnixMilli(req.SendTime),
		UserID:     req.UserID,
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
			o.agentSendConfigMessage(ctx, req.ConversationID, "", func(agent *model.Agent) *model.AgentMessage {
				return agent.TimeoutMsg
			})
		} else {
			o.agentSendConfigMessage(ctx, req.ConversationID, "", func(agent *model.Agent) *model.AgentMessage {
				return agent.EndMsg
			})
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

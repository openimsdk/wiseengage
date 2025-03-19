package convert

import (
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func AgentMessagePb2Db(msg *customerservice.AgentMessage) *model.AgentMessage {
	if msg == nil {
		return nil
	}
	return &model.AgentMessage{
		ContentType: msg.ContentType,
		Content:     msg.Content,
	}
}

func AgentMessageDb2Pb(msg *model.AgentMessage) *customerservice.AgentMessage {
	if msg == nil {
		return nil
	}
	return &customerservice.AgentMessage{
		ContentType: msg.ContentType,
		Content:     msg.Content,
	}
}

func AgentDb2Pb(agent *model.Agent) *customerservice.AgentInfo {
	return &customerservice.AgentInfo{
		UserID:     agent.UserID,
		Nickname:   agent.Nickname,
		FaceURL:    agent.FaceURL,
		AgentType:  agent.Type,
		Status:     agent.Status,
		StartMsg:   AgentMessageDb2Pb(agent.StartMsg),
		EndMsg:     AgentMessageDb2Pb(agent.EndMsg),
		TimeoutMsg: AgentMessageDb2Pb(agent.TimeoutMsg),
		CreateTime: agent.CreateTime.UnixMilli(),
	}
}

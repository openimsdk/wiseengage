package customerservice

import (
	"context"
	"time"

	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/wiseengage/v1/pkg/common/convert"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (o *customerService) CreateAgent(ctx context.Context, req *customerservice.CreateAgentReq) (*customerservice.CreateAgentResp, error) {
	imToken, err := o.imApi.ImAdminTokenWithDefaultAdmin(ctx)
	if err != nil {
		return nil, err
	}
	ctx = mctx.WithApiToken(ctx, imToken)
	if req.Agent.UserID == "" {
		users, err := o.imApi.GetUsers(ctx, []string{req.Agent.UserID})
		if err != nil {
			return nil, err
		}
		if len(users) > 0 {
			return nil, errs.ErrDuplicateKey.WrapMsg("agent userID already exists")
		}
	} else {
		randUserIDs := make([]string, 5)
		for i := range randUserIDs {
			randUserIDs[i] = "bot_" + genID(10)
		}
		users, err := o.imApi.GetUsers(ctx, []string{req.Agent.UserID})
		if err != nil {
			return nil, err
		}
		if len(users) == len(randUserIDs) {
			return nil, errs.ErrDuplicateKey.WrapMsg("gen agent userID already exists")
		}
		for _, user := range users {
			if datautil.Contain(user.UserID, randUserIDs...) {
				continue
			}
			req.Agent.UserID = user.UserID
			break
		}
	}
	agent := &model.Agent{
		UserID:     req.Agent.UserID,
		Nickname:   req.Agent.Nickname,
		FaceURL:    req.Agent.FaceURL,
		Type:       req.Agent.AgentType,
		Status:     req.Agent.Status,
		StartMsg:   convert.AgentMessagePb2Db(req.Agent.StartMsg),
		EndMsg:     convert.AgentMessagePb2Db(req.Agent.EndMsg),
		TimeoutMsg: convert.AgentMessagePb2Db(req.Agent.TimeoutMsg),
		CreateTime: time.Now(),
	}
	userInfo := &sdkws.UserInfo{
		UserID:     req.Agent.UserID,
		Nickname:   req.Agent.Nickname,
		FaceURL:    req.Agent.FaceURL,
		CreateTime: agent.CreateTime.UnixMilli(),
	}
	if err := o.imApi.RegisterUser(ctx, []*sdkws.UserInfo{userInfo}); err != nil {
		return nil, err
	}
	if err := o.db.CreatAgent(ctx, agent); err != nil {
		return nil, err
	}
	return &customerservice.CreateAgentResp{}, nil
}

func (o *customerService) UpdateAgent(ctx context.Context, req *customerservice.UpdateAgentReq) (*customerservice.UpdateAgentResp, error) {
	if _, err := o.db.TakeAgent(ctx, req.UserID); err != nil {
		return nil, err
	}
	update := UpdateAgent(req)
	if len(update) == 0 {
		return nil, errs.ErrArgs.WrapMsg("update data empty")
	}
	if err := o.db.UpdateAgent(ctx, req.UserID, update); err != nil {
		return nil, err
	}
	return &customerservice.UpdateAgentResp{}, nil
}

func (o *customerService) PageFindAgent(ctx context.Context, req *customerservice.PageFindAgentReq) (*customerservice.PageFindAgentResp, error) {
	total, res, err := o.db.PageAgent(ctx, req.AgentTypes, req.Status, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &customerservice.PageFindAgentResp{
		Total:  total,
		Agents: datautil.Slice(res, convert.AgentDb2Pb),
	}, nil
}

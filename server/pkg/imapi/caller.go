package imapi

import (
	"context"
	"sync"
	"time"

	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/tools/utils/jsonutil"

	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/constant"
	constantpb "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
)

type CallerInterface interface {
	ImAdminTokenWithDefaultAdmin(ctx context.Context) (string, error)
	ImportFriend(ctx context.Context, ownerUserID string, friendUserID []string) error
	GetUserToken(ctx context.Context, userID string, platform int32) (string, error)
	GetAdminTokenCache(ctx context.Context, userID string) (string, error)
	InviteToGroup(ctx context.Context, userID string, groupIDs []string) error
	UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error
	ForceOffLine(ctx context.Context, userID string) error
	RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error
	FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error)
	UserRegisterCount(ctx context.Context, start int64, end int64) (map[string]int64, int64, error)
	FriendUserIDs(ctx context.Context, userID string) ([]string, error)
	AccountCheckSingle(ctx context.Context, userID string) (bool, error)
	FindGroupMemberUserIDs(ctx context.Context, groupID string) ([]string, error)
	FindFriendUserIDs(ctx context.Context, userID string) ([]string, error)
	SendBusinessNotification(ctx context.Context, key string, data any, sendUserID string, recvUserID string) error
	SendMsg(ctx context.Context, req *SendMsgReq) (*SendMsgResp, error)
	GetUsers(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error)
	GetFullJoinGroupIDs(ctx context.Context, userID string) ([]string, error)
	CreateGroup(ctx context.Context, req *group.CreateGroupReq) (*sdkws.GroupInfo, error)
	GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error)
	UserRegister(ctx context.Context, user []*sdkws.UserInfo) ([]*sdkws.UserInfo, error)
}

type authToken struct {
	token   string
	timeout time.Time
}

type Caller struct {
	CallerInterface
	imApi           string
	imSecret        string
	defaultIMUserID string
	tokenCache      map[string]*authToken
	lock            sync.RWMutex
}

func (c *Caller) UserRegister(ctx context.Context, user []*sdkws.UserInfo) ([]*sdkws.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Caller) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Caller) GetFullJoinGroupIDs(ctx context.Context, userID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Caller) CreateGroup(ctx context.Context, req *group.CreateGroupReq) (*sdkws.GroupInfo, error) {
	//TODO implement me
	panic("implement me")
}

func New(imApi string, imSecret string, defaultIMUserID string) CallerInterface {
	return &Caller{
		imApi:           imApi,
		imSecret:        imSecret,
		defaultIMUserID: defaultIMUserID,
		tokenCache:      make(map[string]*authToken),
		lock:            sync.RWMutex{},
	}
}

func (c *Caller) ImportFriend(ctx context.Context, ownerUserID string, friendUserIDs []string) error {
	if len(friendUserIDs) == 0 {
		return nil
	}
	_, err := importFriend.Call(ctx, c.imApi, &relation.ImportFriendReq{
		OwnerUserID:   ownerUserID,
		FriendUserIDs: friendUserIDs,
	})
	return err
}

func (c *Caller) ImAdminTokenWithDefaultAdmin(ctx context.Context) (string, error) {
	return c.GetAdminTokenCache(ctx, c.defaultIMUserID)
}

func (c *Caller) GetAdminTokenCache(ctx context.Context, userID string) (string, error) {
	c.lock.RLock()
	t, ok := c.tokenCache[userID]
	c.lock.RUnlock()
	if !ok || t.timeout.Before(time.Now()) {
		c.lock.Lock()
		t, ok = c.tokenCache[userID]
		if !ok || t.timeout.Before(time.Now()) {
			token, err := c.getAdminTokenServer(ctx, userID)
			if err != nil {
				log.ZError(ctx, "get im admin token", err, "userID", userID)
				return "", err
			}
			log.ZDebug(ctx, "get im admin token", "userID", userID)
			t = &authToken{token: token, timeout: time.Now().Add(time.Minute * 5)}
			c.tokenCache[userID] = t
		}
		c.lock.Unlock()
	}
	return t.token, nil
}

func (c *Caller) getAdminTokenServer(ctx context.Context, userID string) (string, error) {
	resp, err := getAdminToken.Call(ctx, c.imApi, &auth.GetAdminTokenReq{
		Secret: c.imSecret,
		UserID: userID,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *Caller) GetUserToken(ctx context.Context, userID string, platformID int32) (string, error) {
	resp, err := getuserToken.Call(ctx, c.imApi, &auth.GetUserTokenReq{
		PlatformID: platformID,
		UserID:     userID,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *Caller) InviteToGroup(ctx context.Context, userID string, groupIDs []string) error {
	for _, groupID := range groupIDs {
		_, _ = inviteToGroup.Call(ctx, c.imApi, &group.InviteUserToGroupReq{
			GroupID:        groupID,
			Reason:         "",
			InvitedUserIDs: []string{userID},
		})
	}
	return nil
}

func (c *Caller) UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error {
	_, err := updateUserInfo.Call(ctx, c.imApi, &user.UpdateUserInfoReq{UserInfo: &sdkws.UserInfo{
		UserID:   userID,
		Nickname: nickName,
		FaceURL:  faceURL,
	}})
	return err
}

func (c *Caller) RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error {
	_, err := registerUser.Call(ctx, c.imApi, &user.UserRegisterReq{
		Users: users,
	})
	return err
}

func (c *Caller) ForceOffLine(ctx context.Context, userID string) error {
	for id := range constantpb.PlatformID2Name {
		_, _ = forceOffLine.Call(ctx, c.imApi, &auth.ForceLogoutReq{
			PlatformID: int32(id),
			UserID:     userID,
		})
	}
	return nil
}

func (c *Caller) FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	resp, err := getGroupsInfo.Call(ctx, c.imApi, &group.GetGroupsInfoReq{
		GroupIDs: groupIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}

func (c *Caller) UserRegisterCount(ctx context.Context, start int64, end int64) (map[string]int64, int64, error) {
	resp, err := registerUserCount.Call(ctx, c.imApi, &user.UserRegisterCountReq{
		Start: start,
		End:   end,
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.Count, resp.Total, nil
}

func (c *Caller) FriendUserIDs(ctx context.Context, userID string) ([]string, error) {
	resp, err := friendUserIDs.Call(ctx, c.imApi, &relation.GetFriendIDsReq{UserID: userID})
	if err != nil {
		return nil, err
	}
	return resp.FriendIDs, nil
}

// return true when isUserNotExist.
func (c *Caller) AccountCheckSingle(ctx context.Context, userID string) (bool, error) {
	resp, err := accountCheck.Call(ctx, c.imApi, &user.AccountCheckReq{CheckUserIDs: []string{userID}})
	if err != nil {
		return false, err
	}
	if resp.Results[0].AccountStatus == constant.Registered {
		return false, eerrs.ErrAccountAlreadyRegister.Wrap()
	}
	return true, nil
}

func (c *Caller) FindGroupMemberUserIDs(ctx context.Context, groupID string) ([]string, error) {
	resp, err := getGroupMemberUserIDs.Call(ctx, c.imApi, &group.GetGroupMemberUserIDsReq{GroupID: groupID})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (c *Caller) FindFriendUserIDs(ctx context.Context, userID string) ([]string, error) {
	resp, err := getFriendIDs.Call(ctx, c.imApi, &relation.GetFriendIDsReq{UserID: userID})
	if err != nil {
		return nil, err
	}
	return resp.FriendIDs, nil
}

func (c *Caller) SendBusinessNotification(ctx context.Context, key string, data any, sendUserID string, recvUserID string) error {
	resp, err := sendBusinessNotification.Call(ctx, c.imApi, &SendBusinessNotificationReq{
		Key: key,
		Data: jsonutil.StructToJsonString(&struct {
			Body any `json:"body"`
		}{Body: data}),
		SendUserID: sendUserID,
		RecvUserID: recvUserID,
	})
	if err != nil {
		log.ZError(ctx, "sendBusinessNotification failed", err, "key", key, "data", data, "sendUserID", sendUserID, "recvUserID", recvUserID)
		return err
	}
	log.ZDebug(ctx, "sendBusinessNotification success", "resp", resp, "key", key, "data", data, "sendUserID", sendUserID, "recvUserID", recvUserID)
	return err
}

func (c *Caller) SendMsg(ctx context.Context, req *SendMsgReq) (*SendMsgResp, error) {
	return sendMsg.Call(ctx, c.imApi, req)
}

func (c *Caller) GetUsers(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	return nil, nil
}

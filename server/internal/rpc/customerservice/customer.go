package customerservice

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/wiseengage/v1/pkg/common/constant"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/controller"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (o *customerService) RegisterCustomer(ctx context.Context, req *customerservice.RegisterCustomerReq) (*customerservice.RegisterCustomerResp, error) {
	now := time.Now()
	imToken, err := o.imApi.ImAdminTokenWithDefaultAdmin(ctx)
	if err != nil {
		return nil, err
	}
	ctx = mctx.WithApiToken(ctx, imToken)
	if req.UserID != "" {
		users, err := o.imApi.GetUsers(ctx, []string{req.UserID})
		if err != nil {
			return nil, err
		}
		if len(users) > 0 {
			return nil, errs.ErrDuplicateKey.WrapMsg("customer userID already exists")
		}
		req.UserID += constant.CustomerUserIDPrefix
	} else {
		randUserIDs := make([]string, 5)
		for i := range randUserIDs {
			randUserIDs[i] = constant.CustomerUserIDPrefix + genID(10)
		}
		users, err := o.imApi.GetUsers(ctx, []string{req.UserID})
		if err != nil {
			return nil, err
		}
		if len(users) == len(randUserIDs) {
			return nil, errs.ErrDuplicateKey.WrapMsg("gen customer userID already exists, please try again")
		}
		for _, user := range users {
			if datautil.Contain(user.UserID, randUserIDs...) {
				continue
			}
			req.UserID = user.UserID
			break
		}
	}
	_, err = o.imApi.UserRegister(ctx, []*sdkws.UserInfo{
		{
			UserID:     req.UserID,
			Nickname:   req.NickName,
			FaceURL:    req.FaceURL,
			Ex:         req.Ex,
			CreateTime: now.UnixMilli(),
		},
	})
	if err != nil {
		return nil, errs.WrapMsg(err, "im register err")
	}

	err = o.db.CustomerCreate(ctx, &model.Customer{
		UserID:     req.UserID,
		NickName:   req.NickName,
		FaceURL:    req.FaceURL,
		Ex:         req.Ex,
		CreateTime: now,
	})
	if err != nil {
		return nil, err
	}

	return &customerservice.RegisterCustomerResp{UserID: req.UserID}, nil
}

func (o *customerService) genCustomerUserID(ctx context.Context) (string, error) {
	const l = 10
	for i := 0; i < 20; i++ {
		userID := "cu" + genID(l)
		_, err := o.db.CustomerTake(ctx, userID)
		if err == nil {
			continue
		} else if controller.IsNotFound(err) {
			return userID, nil
		} else {
			return "", err
		}
	}
	return "", errs.ErrInternalServer.WrapMsg("gen user id failed")
}

func genID(l int) string {
	data := make([]byte, l)
	_, _ = rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}

package customerservice

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/wiseengage/v1/pkg/common/servererrs"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/controller"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (o *customerService) RegisterCustomer(ctx context.Context, req *customerservice.RegisterCustomerReq) (*customerservice.RegisterCustomerResp, error) {
	now := time.Now()

	imReg := false
	if req.UserID != "" {
		exist, err := o.db.CustomerExist(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, servererrs.ErrRegisteredAlready.Wrap()
		}
		u, err := o.userCli.GetUserInfo(ctx, req.UserID)
		if err != nil {
			if !errors.Is(err, errs.ErrRecordNotFound) {
				return nil, err
			}
			if req.UserID == "" {
				req.UserID, err = o.genUserID(ctx)
				if err != nil {
					return nil, err
				}
			}
		} else {
			imReg = true
			req.UserID = u.UserID
		}
	}

	if !imReg {
		_, err := o.userCli.UserRegister(ctx, &user.UserRegisterReq{Users: []*sdkws.UserInfo{
			{
				UserID:     req.UserID,
				Nickname:   req.NickName,
				FaceURL:    req.FaceURL,
				Ex:         req.Ex,
				CreateTime: now.UnixMilli(),
			},
		}})
		if err != nil {
			return nil, errs.WrapMsg(err, "im register err")
		}
	}

	err := o.db.CustomerCreate(ctx, &model.Customer{
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

func (o *customerService) genUserID(ctx context.Context) (string, error) {
	const l = 10
	for i := 0; i < 20; i++ {
		userID := genID(l)
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
	return "cu" + string(data)
}

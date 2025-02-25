package convert

import (
	"time"

	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/common"
)

func CustomerDB2Pb(customer *model.Customer) *common.Customer {
	return &common.Customer{
		UserID:     customer.UserID,
		NickName:   customer.NickName,
		FaceURL:    customer.FaceURL,
		Ex:         customer.Ex,
		CreateTime: customer.CreateTime.UnixMilli(),
	}
}

func CustomersDB2Pb(customers []*model.Customer) []*common.Customer {
	return datautil.Slice(customers, CustomerDB2Pb)
}

func CustomerPb2DB(customer *common.Customer) *model.Customer {
	return &model.Customer{
		UserID:     customer.UserID,
		NickName:   customer.NickName,
		FaceURL:    customer.FaceURL,
		Ex:         customer.Ex,
		CreateTime: time.UnixMilli(customer.CreateTime),
	}
}

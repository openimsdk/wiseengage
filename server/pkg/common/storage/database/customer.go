package database

import (
	"context"

	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
)

type Customer interface {
	Create(ctx context.Context, customers ...*model.Customer) error
	Take(ctx context.Context, userID string) (*model.Customer, error)
	Find(ctx context.Context, userIDs []string) ([]*model.Customer, error)
	Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*model.Customer, err error)
	UpdateByMap(ctx context.Context, userID string, data map[string]any) error
	Delete(ctx context.Context, userIDs []string) error
}

package database

import (
	"context"

	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
)

type Customer interface {
	Create(ctx context.Context, customers ...*model.Customer) error
	Take(ctx context.Context, userID string) (*model.Customer, error)
	Find(ctx context.Context, userIDs []string) ([]*model.Customer, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	Delete(ctx context.Context, userIDs []string) error
}

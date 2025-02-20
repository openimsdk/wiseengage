package database

import (
	"context"

	"wiseengage/server/pkg/common/storage/model"
)

type Customer interface {
	Create(ctx context.Context, jobs ...*model.Customer) error
	Take(ctx context.Context, userID string) (*model.Customer, error)
	Find(ctx context.Context, userIDs []string) ([]*model.Customer, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	Delete(ctx context.Context, userIDs []string) error
}

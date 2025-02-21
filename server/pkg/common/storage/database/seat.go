package database

import (
	"context"

	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
)

type Seat interface {
	Create(ctx context.Context, jobs ...*model.Seat) error
	Take(ctx context.Context, userID string) (*model.Seat, error)
	Find(ctx context.Context, userIDs []string) ([]*model.Seat, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	UpdatePassword(ctx context.Context, userId string, password string) error
	Delete(ctx context.Context, userIDs []string) error
}

package database

import (
	"context"

	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
)

type Agent interface {
	Create(ctx context.Context, jobs ...*model.Agent) error
	Take(ctx context.Context, userID string) (*model.Agent, error)
	Find(ctx context.Context, userIDs []string) ([]*model.Agent, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	Delete(ctx context.Context, userIDs []string) error
}

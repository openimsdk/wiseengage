package database

import (
	"context"

	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
)

type Agent interface {
	Create(ctx context.Context, jobs ...*model.Agent) error
	Take(ctx context.Context, userID string) (*model.Agent, error)
	Find(ctx context.Context, userIDs []string) ([]*model.Agent, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	Delete(ctx context.Context, userIDs []string) error
	Page(ctx context.Context, types []string, status []string, pagination pagination.Pagination) (int64, []*model.Agent, error)
	FindType(ctx context.Context, agentType string, status []string) ([]*model.Agent, error)
}

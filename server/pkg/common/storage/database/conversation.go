package database

import (
	"context"
	"time"

	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
)

type Conversation interface {
	Take(ctx context.Context, conversationID string) (*model.Conversation, error)
	Create(ctx context.Context, conversation *model.Conversation) error
	TakeByUserID(ctx context.Context, userID string) (*model.Conversation, error)
	UpdateLastMsg(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage) error
	SetStatusOpen(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error)
	SetStatusClosed(ctx context.Context, userID string, conversationID string, version int, cause string) (bool, error)
	SetStatusTimeoutClosed(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage, cause string) (bool, error)
	SetRole(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error)
	FindTimeout(ctx context.Context, deadline time.Time, limit int) ([]*model.Conversation, error)
	Find(ctx context.Context, conversationIDs []string) ([]*model.Conversation, error)
}

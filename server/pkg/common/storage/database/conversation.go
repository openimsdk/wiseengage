package database

import (
	"context"

	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
)

type Conversation interface {
	Create(ctx context.Context, conversation *model.Conversation) error
	TakeByUserID(ctx context.Context, userID string) (*model.Conversation, error)
	UpdateStatus(ctx context.Context, userID string, conversationID string, status int, role string, version int) (bool, error)
	UpdateLastMsg(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage) error
}

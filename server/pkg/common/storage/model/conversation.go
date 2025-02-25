package model

import (
	"time"
)

type LastMessage struct {
	Seq        int64     `bson:"seq"`
	SendTime   time.Time `bson:"send_time"`
	UserID     string    `bson:"user_id"`
	UpdateTime time.Time `bson:"update_time"`
}

type Conversation struct {
	UserID         string       `bson:"user_id"`
	ConversationID string       `bson:"conversation_id"`
	CreateTime     time.Time    `bson:"create_time"`
	LastMsg        *LastMessage `bson:"last_msg"`
	Status         int          `bson:"status"`
	Role           string       `bson:"role"`
	Cause          string       `bson:"cause"`
	Version        int          `bson:"version"`
	UpdateTime     time.Time    `bson:"update_time"`
}

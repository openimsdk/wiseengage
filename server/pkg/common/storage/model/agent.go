package model

import "time"

type Agent struct {
	UserID     string        `bson:"user_id"`
	Nickname   string        `bson:"nickname"`
	FaceURL    string        `bson:"face_url"`
	Type       string        `bson:"type"`
	Status     int32         `bson:"status"`
	StartMsg   *AgentMessage `bson:"start_msg"`
	EndMsg     *AgentMessage `bson:"end_msg"`
	TimeoutMsg *AgentMessage `bson:"timeout_msg"`
	CreateTime time.Time     `bson:"create_time"`
}

type AgentMessage struct {
	ContentType int32  `bson:"content_type"`
	Content     string `bson:"content"`
}

package model

import (
	"time"
)

type Seat struct {
	UserID         string    `bson:"user_id"`
	Name           string    `bson:"name"`
	ShowName       string    `bson:"show_name"`
	Account        string    `bson:"account"`
	Password       string    `bson:"password"`
	CreateTime     time.Time `bson:"create_time"`
	ChangeTime     time.Time `bson:"change_time"`
	OperatorUserID string    `bson:"operator_user_id"`
}

package model

import "time"

type Customer struct {
	UserID     string    `bson:"user_id"`
	NickName   string    `bson:"nick_name"`
	FaceURL    string    `bson:"face_url"`
	Ex         string    `bson:"ex"`
	CreateTime time.Time `bson:"create_time"`
}

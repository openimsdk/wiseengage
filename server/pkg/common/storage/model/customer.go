package model

type Customer struct {
	UserID   string `bson:"user_id"`
	NickName string `bson:"nick_name"`
	Key      string `bson:"key"`
	FaceURL  string `bson:"face_url"`
}

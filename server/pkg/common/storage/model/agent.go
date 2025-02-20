package model

type Agent struct {
	UserID string `bson:"user_id"`
	Name   string `bson:"name"`
	URL    string `bson:"url"`
}

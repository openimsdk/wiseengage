package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/database"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewConversation(db *mongo.Database) (database.Conversation, error) {
	coll := db.Collection(database.Prefix + database.ConversationName)
	return &Conversation{coll: coll}, nil
}

type Conversation struct {
	coll *mongo.Collection
}

func (c *Conversation) Create(ctx context.Context, conversation *model.Conversation) error {
	conversation.UpdateTime = conversation.CreateTime
	return mongoutil.InsertOne(ctx, c.coll, conversation)
}

func (c *Conversation) Take(ctx context.Context, userID string, conversationID string) (*model.Conversation, error) {
	return mongoutil.FindOne[*model.Conversation](ctx, c.coll, bson.M{"user_id": userID, "conversation_id": conversationID})
}

func (c *Conversation) TakeByUserID(ctx context.Context, userID string) (*model.Conversation, error) {
	return mongoutil.FindOne[*model.Conversation](ctx, c.coll, bson.M{"user_id": userID})
}

func (c *Conversation) UpdateStatus(ctx context.Context, userID string, conversationID string, status int, role string, version int) (bool, error) {
	filter := bson.M{
		"user_id":         userID,
		"conversation_id": conversationID,
	}
	if version > 0 {
		filter["version"] = version
	}
	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"update_time": time.Now(),
		},
		"$inc": bson.M{"version": 1},
	}
	result, err := mongoutil.UpdateOneResult(ctx, c.coll, filter, update)
	if err != nil {
		return false, err
	}
	return result.UpsertedCount > 0, nil
}

func (c *Conversation) UpdateLastMsg(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage) error {
	filter := bson.M{
		"user_id":         userID,
		"conversation_id": conversationID,
		"$or": []bson.M{
			{"last_msg": nil},
			{"last_msg.seq": bson.M{"$lt": lastMsg.Seq}},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"last_msg":    lastMsg,
			"update_time": time.Now(),
		},
	}
	return mongoutil.UpdateOne(ctx, c.coll, filter, update, false)
}

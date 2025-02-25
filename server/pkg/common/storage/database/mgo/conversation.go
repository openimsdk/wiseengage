package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/wiseengage/v1/pkg/common/constant"
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

func (c *Conversation) getFilter(userID string, conversationID string, version int) bson.M {
	filter := bson.M{
		"user_id":         userID,
		"conversation_id": conversationID,
	}
	if version > 0 {
		filter["version"] = version
	}
	return filter
}

func (c *Conversation) UpdateLastMsg(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage) error {
	filter := c.getFilter(userID, conversationID, -1)
	filter["$or"] = []bson.M{
		{"last_msg": nil},
		{"last_msg.seq": bson.M{"$lt": lastMsg.Seq}},
	}

	update := bson.M{
		"last_msg": lastMsg,
	}
	_, err := c.update(ctx, filter, false, false, update)
	return err
}

func (c *Conversation) update(ctx context.Context, filter bson.M, updateTime bool, version bool, data map[string]any) (bool, error) {
	if len(data) > 0 {
		delete(data, "update_time")
	}
	if len(data) == 0 {
		return false, nil
	}
	if updateTime {
		data["update_time"] = time.Now()
	}
	update := bson.M{
		"$set": data,
	}
	if version {
		update["$set"] = bson.M{"version": 1}
	}
	result, err := mongoutil.UpdateOneResult(ctx, c.coll, filter, update)
	if err != nil {
		return false, err
	}
	return result.UpsertedCount > 0, nil
}

func (c *Conversation) SetStatusOpen(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error) {
	set := bson.M{
		"status":      constant.ConversationStatusClosed,
		"update_time": time.Now(),
	}
	if role != "" {
		set["role"] = role
	}
	return c.update(ctx, c.getFilter(userID, conversationID, version), true, true, set)
}

func (c *Conversation) SetStatusClosed(ctx context.Context, userID string, conversationID string, version int, cause string) (bool, error) {
	set := bson.M{
		"status":      constant.ConversationStatusClosed,
		"cause":       cause,
		"update_time": time.Now(),
	}
	return c.update(ctx, c.getFilter(userID, conversationID, version), true, true, set)
}

func (c *Conversation) SetRole(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error) {
	filter := c.getFilter(userID, conversationID, version)
	filter["status"] = constant.ConversationStatusOpen
	filter["role"] = bson.M{"$ne": role}
	set := bson.M{
		"role":        role,
		"update_time": time.Now(),
	}
	return c.update(ctx, filter, true, true, set)
}

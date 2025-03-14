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
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (c *Conversation) Take(ctx context.Context, conversationID string) (*model.Conversation, error) {
	return mongoutil.FindOne[*model.Conversation](ctx, c.coll, bson.M{"conversation_id": conversationID})
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
		{"last_msg.send_time": bson.M{"$lt": lastMsg.SendTime}},
	}
	update := bson.M{
		"last_msg": lastMsg,
	}
	_, err := c.update(ctx, filter, false, update)
	return err
}

func (c *Conversation) update(ctx context.Context, filter bson.M, version bool, data map[string]any) (bool, error) {
	if len(data) == 0 {
		return false, nil
	}
	data["update_time"] = time.Now()
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
		"status":      constant.ConversationStatusOpen,
		"update_time": time.Now(),
	}
	if role != "" {
		set["role"] = role
	}
	return c.update(ctx, c.getFilter(userID, conversationID, version), true, set)
}

func (c *Conversation) SetStatusClosed(ctx context.Context, userID string, conversationID string, version int, cause string) (bool, error) {
	set := bson.M{
		"status":      constant.ConversationStatusClosed,
		"cause":       cause,
		"update_time": time.Now(),
	}
	return c.update(ctx, c.getFilter(userID, conversationID, version), true, set)
}

func (c *Conversation) SetStatusTimeoutClosed(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage, cause string) (bool, error) {
	filter := c.getFilter(userID, conversationID, -1)
	if lastMsg == nil {
		filter["last_msg"] = nil
	} else {
		filter["last_msg.msg_id"] = lastMsg.MsgID
	}
	set := bson.M{
		"status":      constant.ConversationStatusClosed,
		"cause":       cause,
		"update_time": time.Now(),
	}
	return c.update(ctx, c.getFilter(userID, conversationID, -1), true, set)
}

func (c *Conversation) SetRole(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error) {
	filter := c.getFilter(userID, conversationID, version)
	filter["status"] = constant.ConversationStatusOpen
	filter["role"] = bson.M{"$ne": role}
	set := bson.M{
		"role":        role,
		"update_time": time.Now(),
	}
	return c.update(ctx, filter, true, set)
}

func (c *Conversation) FindTimeout(ctx context.Context, deadline time.Time, limit int) ([]*model.Conversation, error) {
	filter := bson.M{
		"status": constant.ConversationStatusOpen,
		"update_time": bson.M{
			"$lt": deadline,
		},
	}
	opts := options.Find().SetLimit(int64(limit))
	return mongoutil.Find[*model.Conversation](ctx, c.coll, filter, opts)
}

func (c *Conversation) Find(ctx context.Context, conversationIDs []string) ([]*model.Conversation, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	filter := bson.M{
		"conversation_id": bson.M{"$in": conversationIDs},
	}
	return mongoutil.Find[*model.Conversation](ctx, c.coll, filter)
}

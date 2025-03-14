package mgo

import (
	"context"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/database"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAgent(db *mongo.Database) (database.Agent, error) {
	coll := db.Collection("agent")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Agent{coll: coll}, nil
}

type Agent struct {
	coll *mongo.Collection
}

func (o *Agent) Create(ctx context.Context, accounts ...*model.Agent) error {
	return mongoutil.InsertMany(ctx, o.coll, accounts)
}

func (o *Agent) Take(ctx context.Context, userId string) (*model.Agent, error) {
	return mongoutil.FindOne[*model.Agent](ctx, o.coll, bson.M{"user_id": userId})
}

func (o *Agent) Find(ctx context.Context, userIDs []string) ([]*model.Agent, error) {
	return mongoutil.Find[*model.Agent](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Agent) FindType(ctx context.Context, agentType string, status []string) ([]*model.Agent, error) {
	filter := bson.M{
		"type": agentType,
	}
	if len(status) > 0 {
		filter["status"] = bson.M{"$in": status}
	}
	return mongoutil.Find[*model.Agent](ctx, o.coll, filter)
}

func (o *Agent) Update(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (o *Agent) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Agent) Page(ctx context.Context, types []string, status []string, pagination pagination.Pagination) (int64, []*model.Agent, error) {
	filter := bson.M{}
	if len(types) > 0 {
		filter["type"] = bson.M{"$in": types}
	}
	if len(status) > 0 {
		filter["status"] = bson.M{"$in": status}
	}
	return mongoutil.FindPage[*model.Agent](ctx, o.coll, filter, pagination)
}

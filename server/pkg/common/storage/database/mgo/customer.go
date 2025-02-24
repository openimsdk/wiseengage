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

func NewCustomer(db *mongo.Database) (database.Customer, error) {
	coll := db.Collection("account")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Customer{coll: coll}, nil
}

type Customer struct {
	coll *mongo.Collection
}

func (c *Customer) Create(ctx context.Context, accounts ...*model.Customer) error {
	return mongoutil.InsertMany(ctx, c.coll, accounts)
}

func (c *Customer) Take(ctx context.Context, userId string) (*model.Customer, error) {
	return mongoutil.FindOne[*model.Customer](ctx, c.coll, bson.M{"user_id": userId})
}

func (c *Customer) Find(ctx context.Context, userIDs []string) ([]*model.Customer, error) {
	return mongoutil.Find[*model.Customer](ctx, c.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (c *Customer) Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*model.Customer, err error) {
	return mongoutil.FindPage[*model.Customer](ctx, c.coll, bson.M{}, pagination)
}

func (c *Customer) UpdateByMap(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, c.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (c *Customer) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, c.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

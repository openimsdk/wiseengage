package mgo

import (
	"context"

	"wiseengage/server/pkg/common/storage/database"
	"wiseengage/server/pkg/common/storage/model"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
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

func (o *Customer) Create(ctx context.Context, accounts ...*model.Customer) error {
	return mongoutil.InsertMany(ctx, o.coll, accounts)
}

func (o *Customer) Take(ctx context.Context, userId string) (*model.Customer, error) {
	return mongoutil.FindOne[*model.Customer](ctx, o.coll, bson.M{"user_id": userId})
}

func (o *Customer) Find(ctx context.Context, userIDs []string) ([]*model.Customer, error) {
	return mongoutil.Find[*model.Customer](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Customer) Update(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (o *Customer) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

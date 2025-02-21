package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/database"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewSeat(db *mongo.Database) (database.Seat, error) {
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
	return &Seat{coll: coll}, nil
}

type Seat struct {
	coll *mongo.Collection
}

func (o *Seat) Create(ctx context.Context, accounts ...*model.Seat) error {
	return mongoutil.InsertMany(ctx, o.coll, accounts)
}

func (o *Seat) Take(ctx context.Context, userId string) (*model.Seat, error) {
	return mongoutil.FindOne[*model.Seat](ctx, o.coll, bson.M{"user_id": userId})
}

func (o *Seat) Find(ctx context.Context, userIDs []string) ([]*model.Seat, error) {
	return mongoutil.Find[*model.Seat](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Seat) Update(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (o *Seat) UpdatePassword(ctx context.Context, userId string, password string) error {
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"password": password, "change_time": time.Now()}}, false)
}

func (o *Seat) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

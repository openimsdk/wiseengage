package controller

import (
	"context"
	"errors"

	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/database"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsNotFound(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

func IsDuplicateKeyError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

type CustomerDatabase interface {
	TakeConversationByUserID(ctx context.Context, userID string) (*model.Conversation, error)
	CreateConversation(ctx context.Context, conversation *model.Conversation) error
	UpdateConversationLastMsg(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage) error
	UpdateConversationStatusOpen(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error)
	UpdateConversationStatusClosed(ctx context.Context, userID string, conversationID string, version int, cause string) (bool, error)
	UpdateConversationRole(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error)

	CustomerCreate(ctx context.Context, customers ...*model.Customer) (err error)
	CustomerUpdateByMap(ctx context.Context, customerID string, args map[string]any) (err error)
	CustomerPage(ctx context.Context, pagination pagination.Pagination) (count int64, customers []*model.Customer, err error)
	CustomerTake(ctx context.Context, customerID string) (customers *model.Customer, err error)
	CustomerFind(ctx context.Context, customerIDs []string) (customers []*model.Customer, err error)
	CustomerExist(ctx context.Context, customerID string) (exist bool, err error)
}

type customerDatabase struct {
	tx             tx.Tx
	customerDB     database.Customer
	conversationDB database.Conversation
}

func NewCustomerDatabase(CustomerDB database.Customer, tx tx.Tx) CustomerDatabase {
	return &customerDatabase{customerDB: CustomerDB, tx: tx}
}

func (c *customerDatabase) TakeConversationByUserID(ctx context.Context, userID string) (*model.Conversation, error) {
	return c.conversationDB.TakeByUserID(ctx, userID)
}

func (c *customerDatabase) CreateConversation(ctx context.Context, conversation *model.Conversation) error {
	return c.conversationDB.Create(ctx, conversation)
}

func (c *customerDatabase) UpdateConversationLastMsg(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage) error {
	return c.conversationDB.UpdateLastMsg(ctx, userID, conversationID, lastMsg)
}
func (c *customerDatabase) UpdateConversationStatusOpen(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error) {
	return c.conversationDB.SetStatusOpen(ctx, userID, conversationID, version, role)
}

func (c *customerDatabase) UpdateConversationStatusClosed(ctx context.Context, userID string, conversationID string, version int, cause string) (bool, error) {
	return c.conversationDB.SetStatusClosed(ctx, userID, conversationID, version, cause)
}

func (c *customerDatabase) UpdateConversationRole(ctx context.Context, userID string, conversationID string, version int, role string) (bool, error) {
	return c.conversationDB.SetRole(ctx, userID, conversationID, version, role)
}

// CustomerCreate Insert multiple external guarantees that the customerID is not repeated and does not exist in the storage.
func (c *customerDatabase) CustomerCreate(ctx context.Context, customers ...*model.Customer) (err error) {
	if err = c.customerDB.Create(ctx, customers...); err != nil {
		return err
	}
	return nil
}

// CustomerUpdateByMap update (zero value) externally guarantees that customerID exists.
func (c *customerDatabase) CustomerUpdateByMap(ctx context.Context, customerID string, args map[string]any) (err error) {
	if err = c.customerDB.UpdateByMap(ctx, customerID, args); err != nil {
		return err
	}
	return nil
}

// CustomerPage Gets, returns no error if not found.
func (c *customerDatabase) CustomerPage(ctx context.Context, pagination pagination.Pagination) (count int64, customers []*model.Customer, err error) {
	return c.customerDB.Page(ctx, pagination)
}

func (c *customerDatabase) CustomerTake(ctx context.Context, customerID string) (customers *model.Customer, err error) {
	return c.customerDB.Take(ctx, customerID)
}

// CustomerFind Get the information of the specified customer. If the customerID is not found, no error will be returned.
func (c *customerDatabase) CustomerFind(ctx context.Context, customerIDs []string) (customers []*model.Customer, err error) {
	return c.customerDB.Find(ctx, customerIDs)
}

func (c *customerDatabase) CustomerExist(ctx context.Context, customerID string) (exist bool, err error) {
	res, err := c.customerDB.Take(ctx, customerID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return false, err
	}
	return res != nil, nil
}

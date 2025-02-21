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
	UpdateConversationStatus(ctx context.Context, userID string, conversationID string, status int, role string, version int) (bool, error)
	UpdateConversationLastMsg(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage) error
}

type customerDatabase struct {
	tx             tx.Tx
	customerDB     database.Customer
	conversationDB database.Conversation
}

func NewCustomerDatabase(CustomerDB database.Customer, tx tx.Tx) CustomerDatabase {
	return &customerDatabase{customerDB: CustomerDB, tx: tx}
}

// Create Insert multiple external guarantees that the customerID is not repeated and does not exist in the storage.
func (u *customerDatabase) Create(ctx context.Context, customers []*model.Customer) (err error) {
	if err = u.customerDB.Create(ctx, customers...); err != nil {
		return err
	}
	return nil
}

// UpdateByMap update (zero value) externally guarantees that customerID exists.
func (u *customerDatabase) UpdateByMap(ctx context.Context, customerID string, args map[string]any) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := u.customerDB.UpdateByMap(ctx, customerID, args); err != nil {
			return err
		}
		return u.cache.DelCustomersInfo(customerID).ChainExecDel(ctx)
	})
}

// Page Gets, returns no error if not found.
func (u *customerDatabase) Page(ctx context.Context, pagination pagination.Pagination) (count int64, customers []*model.Customer, err error) {
	return u.customerDB.Page(ctx, pagination)
}

func (u *customerDatabase) PageFindCustomer(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, customers []*model.Customer, err error) {
	return u.customerDB.PageFindCustomer(ctx, level1, level2, pagination)
}

// Find Get the information of the specified customer. If the customerID is not found, no error will be returned.
func (u *customerDatabase) Find(ctx context.Context, customerIDs []string) (customers []*model.Customer, err error) {
	return u.cache.GetCustomersInfo(ctx, customerIDs)
}

func (u *customerDatabase) TakeConversationByUserID(ctx context.Context, userID string) (*model.Conversation, error) {
	return u.conversationDB.TakeByUserID(ctx, userID)
}

func (u *customerDatabase) CreateConversation(ctx context.Context, conversation *model.Conversation) error {
	return u.conversationDB.Create(ctx, conversation)
}

func (u *customerDatabase) UpdateConversationStatus(ctx context.Context, userID string, conversationID string, status int, role string, version int) (bool, error) {
	return u.conversationDB.UpdateStatus(ctx, userID, conversationID, status, role, version)
}

func (u *customerDatabase) UpdateConversationLastMsg(ctx context.Context, userID string, conversationID string, lastMsg *model.LastMessage) error {
	return u.conversationDB.UpdateLastMsg(ctx, userID, conversationID, lastMsg)
}

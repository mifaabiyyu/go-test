package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mifaabiyyu/go-test.git/controllers"
	"github.com/mifaabiyyu/go-test.git/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionCollection interface wraps MongoDB collection methods
type TransactionCollection interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	Database() *mongo.Database
}

// MockCollection is a mock implementation of TransactionCollection
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockCollection) Database() *mongo.Database {
	args := m.Called()
	return args.Get(0).(*mongo.Database)
}

type MockCursor struct {
	mock.Mock
}

func (m *MockCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockCursor) Decode(val interface{}) error {
	args := m.Called(val)
	return args.Error(0)
}

func TestGetTransactions(t *testing.T) {
	mockCollection := new(MockCollection)
	mockCursor := new(MockCursor)
	mockContext := &gin.Context{}

	tc := &controllers.TransactionController{
		Collection: MockCollection.Mock,
	}

	mockCollection.On("Find", mock.Anything, mock.Anything, mock.AnythingOfType("*options.FindOptions")).Return(mockCursor, nil).Once()
	mockCursor.On("Next", mock.Anything).Return(true).Once()
	mockCursor.On("Next", mock.Anything).Return(false)
	mockCursor.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		transaction := args.Get(0).(*models.Transaction)
		// Initialize transaction fields as needed for your test
		transaction.ID = "123"
	})

	// Assuming your transaction model has been defined in the "models" package
	var transactions []models.Transaction
	mockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
		transactions = args.Get(1).([]models.Transaction)
	})

	tc.GetTransactions(mockContext)

	assert.NotNil(t, transactions)
	assert.Equal(t, 1, len(transactions))
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
	mockContext.AssertExpectations(t)
}

// controllers/transaction_payment_controller.go
package controllers

import (
	"context"
	// "net/http"

	// "github.com/gin-gonic/gin"
	"github.com/mifaabiyyu/go-test.git/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionPaymentController struct {
	Collection *mongo.Collection
}

func NewTransactionPaymentController(db *mongo.Database) *TransactionPaymentController {
	return &TransactionPaymentController{
		Collection: db.Collection("transaction_payments"),
	}
}

func (tpc *TransactionPaymentController) CreateTransactionPayment(payment *models.TransactionPayment) error {
	_, err := tpc.Collection.InsertOne(context.Background(), payment)
	if err != nil {
		return err
	}
	return nil
}

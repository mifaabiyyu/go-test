// models/transaction.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty"  json:"id"`
	CustomerID      primitive.ObjectID   `bson:"customer_id" binding:"required" json:"customer_id"`
	TotalAmount     float64              `bson:"total_amount"  json:"total_amount"`
	TotalQty        float64              `bson:"total_qty"  json:"total_qty"`
	TransactionDate time.Time            `bson:"transaction_date"  json:"transaction_date"`
	Details         []TransactionDetail  `bson:"details"  json:"details"`
	Payments        []TransactionPayment `bson:"payments"  json:"payments"`
}

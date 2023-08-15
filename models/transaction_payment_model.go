package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionPayment struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionID   primitive.ObjectID `bson:"transaction_id" binding:"required" json:"transaction_id"`
	PaymentMethodID primitive.ObjectID `bson:"payment_method_id" binding:"required" json:"payment_method_id"`
	Status          int64              `bson:"status" binding:"required" json:"status"`
	PaidAmount      float64            `bson:"paid_amount" binding:"required" json:"paid_amount"`
	PaymentDate     time.Time          `bson:"payment_date"  json:"payment_date"`
}

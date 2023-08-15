// models/transaction_detail.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TransactionDetail struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	TransactionID primitive.ObjectID `bson:"transaction_id" json:"transaction_id"`
	ProductID     primitive.ObjectID `bson:"product_id" json:"product_id"`
	Quantity      float64            `bson:"quantity" binding:"required" json:"quantity"`
	Subtotal      float64            `bson:"subtotal" binding:"required" json:"subtotal"`
	Price         float64            `bson:"price" binding:"required" json:"price"`
}

// models/payment_method.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentMethod struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	IsActive bool               `bson:"is_active" json:"is_active"`
}

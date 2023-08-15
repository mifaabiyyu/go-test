package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerAddress struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	Street     string             `json:"street"`
	City       string             `json:"city"`
	PostalCode string             `json:"postal_code"`
}

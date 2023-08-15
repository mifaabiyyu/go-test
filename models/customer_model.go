package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" binding:"required"`
	Code  string             `json:"code" binding:"required"`
	Email string             `json:"email" bson:"email,omitempty" binding:"required" unique:"true"`
}


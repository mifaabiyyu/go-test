package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string             `bson:"code" binding:"required" json:"code"`
	Name        string             `bson:"name" binding:"required" json:"name"`
	Price       float64            `bson:"price" binding:"required" json:"price"`
	Description string             `bson:"description" json:"description"`
}

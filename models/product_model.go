package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code 	string          		`bson:"code" json:"code"`
	Name     string             `bson:"name" json:"name"`
	Price    float64            `bson:"price" json:"price"`
	Description string          `bson:"description" json:"description"`
}

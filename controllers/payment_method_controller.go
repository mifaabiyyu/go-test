// controllers/payment_method_controller.go
package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mifaabiyyu/go-test.git/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentMethodController struct {
	Collection *mongo.Collection
}

func NewPaymentMethodController(db *mongo.Database) *PaymentMethodController {
	return &PaymentMethodController{
		Collection: db.Collection("payment_methods"),
	}
}

func (pmc *PaymentMethodController) CreatePaymentMethod(c *gin.Context) {
	var paymentMethod models.PaymentMethod
	if err := c.ShouldBindJSON(&paymentMethod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingPaymentMethod := models.PaymentMethod{}
	err := pmc.Collection.FindOne(context.Background(), bson.M{"name": paymentMethod.Name}).Decode(&existingPaymentMethod)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Name already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Clear the ID field to let MongoDB generate it
	paymentMethod.ID = primitive.NewObjectID()

	_, err = pmc.Collection.InsertOne(context.Background(), paymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentMethod)
}

func (pmc *PaymentMethodController) GetPaymentMethods(c *gin.Context) {
	var paymentMethods []models.PaymentMethod
	cursor, err := pmc.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var paymentMethod models.PaymentMethod
		cursor.Decode(&paymentMethod)
		paymentMethods = append(paymentMethods, paymentMethod)
	}

	c.JSON(http.StatusOK, paymentMethods)
}

func (pmc *PaymentMethodController) GetPaymentMethod(c *gin.Context) {
	id := c.Param("id")
	paymentMethodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID format"})
		return
	}

	var paymentMethod models.PaymentMethod
	err = pmc.Collection.FindOne(context.Background(), bson.M{"_id": paymentMethodID}).Decode(&paymentMethod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment Method not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}

func (pmc *PaymentMethodController) UpdatePaymentMethod(c *gin.Context) {
	id := c.Param("id")
	paymentMethodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID format"})
		return
	}

	var updatedPaymentMethod models.PaymentMethod
	if err := c.ShouldBindJSON(&updatedPaymentMethod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = pmc.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": paymentMethodID},
		bson.M{"$set": updatedPaymentMethod},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPaymentMethod)
}

func (pmc *PaymentMethodController) DeletePaymentMethod(c *gin.Context) {
	id := c.Param("id")
	paymentMethodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID format"})
		return
	}

	// Check if the payment method exists
	var existingPaymentMethod models.PaymentMethod
	err = pmc.Collection.FindOne(context.Background(), bson.M{"_id": paymentMethodID}).Decode(&existingPaymentMethod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Proceed with payment method deletion
	_, err = pmc.Collection.DeleteOne(context.Background(), bson.M{"_id": paymentMethodID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted"})
}

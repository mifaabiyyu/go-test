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

type CustomerAddressController struct {
	Collection *mongo.Collection
}

func NewCustomerAddressController(db *mongo.Database) *CustomerAddressController {
	return &CustomerAddressController{
		Collection: db.Collection("customer_address"),
	}
}

func (cac *CustomerAddressController) CreateCustomerAddress(c *gin.Context) {
	var customerAddress models.CustomerAddress
	if err := c.ShouldBindJSON(&customerAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customerAddress.ID = primitive.NewObjectID()
	_, err := cac.Collection.InsertOne(context.Background(), customerAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customerAddress)
}

func (cac *CustomerAddressController) GetCustomerAddresses(c *gin.Context) {
	var customerAddresses []models.CustomerAddress
	cursor, err := cac.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var customerAddress models.CustomerAddress
		cursor.Decode(&customerAddress)
		customerAddresses = append(customerAddresses, customerAddress)
	}

	c.JSON(http.StatusOK, customerAddresses)
}

func (cac *CustomerAddressController) GetCustomerAddress(c *gin.Context) {
	id := c.Param("id")
	idAddress, err := primitive.ObjectIDFromHex(id)
	var customerAddress models.CustomerAddress
	err = cac.Collection.FindOne(context.Background(), bson.M{"_id": idAddress}).Decode(&customerAddress)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer address not found"})
		return
	}

	c.JSON(http.StatusOK, customerAddress)
}

func (cac *CustomerAddressController) UpdateCustomerAddress(c *gin.Context) {
	id := c.Param("id")
	customerAddressID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer address ID format"})
		return
	}

	var customerAddress models.CustomerAddress
	if err := c.ShouldBindJSON(&customerAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the customer address exists
	existingCustomerAddress := models.CustomerAddress{}
	err = cac.Collection.FindOne(context.Background(), bson.M{"_id": customerAddressID}).Decode(&existingCustomerAddress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = cac.Collection.ReplaceOne(
		context.Background(),
		bson.M{"_id": customerAddressID},
		customerAddress,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	customerAddress.ID = primitive.NewObjectID()

	c.JSON(http.StatusOK, customerAddress)
}

func (cac *CustomerAddressController) DeleteCustomerAddress(c *gin.Context) {
	id := c.Param("id")
	customerAddressID, err := primitive.ObjectIDFromHex(id)

	var existingAddress models.CustomerAddress
	err = cac.Collection.FindOne(context.Background(), bson.M{"_id": customerAddressID}).Decode(&existingAddress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = cac.Collection.DeleteOne(context.Background(), bson.M{"_id": customerAddressID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer address deleted"})
}

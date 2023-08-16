package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mifaabiyyu/go-test.git/models"
)

type CustomerController struct {
	Collection *mongo.Collection
}

type CustomerWithAddresses struct {
	ID        primitive.ObjectID       `bson:"_id,omitempty" json:"id"`
	FirstName string                   `bson:"first_name" json:"first_name"`
	LastName  string                   `bson:"last_name" json:"last_name"`
	Email     string                   `bson:"email" json:"email"`
	Addresses []models.CustomerAddress `bson:"addresses" json:"addresses"`
}

func NewCustomerController(db *mongo.Database) *CustomerController {
	return &CustomerController{
		Collection: db.Collection("customers"),
	}
}

func (cc *CustomerController) CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingCustomer := models.Customer{}
	err := cc.Collection.FindOne(context.Background(), bson.M{"email": customer.Email}).Decode(&existingCustomer)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	customer.ID = primitive.NewObjectID()
	_, err = cc.Collection.InsertOne(context.Background(), customer)
	if err != nil {
		// Check if the error is due to duplicate email
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (cc *CustomerController) GetCustomers(c *gin.Context) {
	pipeline := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "customer_address", // Name of the customer_address collection
				"localField":   "_id",
				"foreignField": "customer_id",
				"as":           "addresses",
			},
		},
	}

	// Perform the aggregation
	cursor, err := cc.Collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var customers []CustomerWithAddresses
	if err := cursor.All(context.Background(), &customers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (cc *CustomerController) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID format"})
		return
	}

	// Aggregate pipeline to join customer with associated addresses
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"_id": customerID},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "customer_address", // Name of the customer_address collection
				"localField":   "_id",
				"foreignField": "customer_id",
				"as":           "addresses",
			},
		},
	}

	// Perform the aggregation
	cursor, err := cc.Collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var customerWithAddresses CustomerWithAddresses
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&customerWithAddresses); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customerWithAddresses)
}

func (cc *CustomerController) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID format"})
		return
	}

	// Check if the customer exists
	existingCustomer := models.Customer{}
	err = cc.Collection.FindOne(context.Background(), bson.M{"_id": customerID}).Decode(&existingCustomer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedCustomer models.Customer
	if err := c.ShouldBindJSON(&updatedCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = cc.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": customerID},
		bson.M{"$set": updatedCustomer},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedCustomer.ID = primitive.NewObjectID()
	c.JSON(http.StatusOK, updatedCustomer)
}

func (cc *CustomerController) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	customerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID format"})
		return
	}
	clientOptions := options.Client().ApplyURI("mongodb+srv://mifaabiyyu:indun1234@cluster0.geg1dqt.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db := client.Database("mydatabase")
	cac := NewTransactionController(cc.Collection.Database(), NewTransactionPaymentController(db))

	var existingCustomer models.Customer
	err = cc.Collection.FindOne(context.Background(), bson.M{"_id": customerID}).Decode(&existingCustomer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transactionCount, err := cac.Collection.CountDocuments(context.Background(), bson.M{"customer_id": customerID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if transactionCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete customer with associated transaction"})
		return
	}

	// Delete associated customer addresses
	_, err = cc.Collection.Database().Collection("customer_addresses").DeleteMany(context.Background(), bson.M{"customer_id": customerID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Proceed with customer deletion
	_, err = cc.Collection.DeleteOne(context.Background(), bson.M{"_id": customerID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer and associated addresses deleted"})
}

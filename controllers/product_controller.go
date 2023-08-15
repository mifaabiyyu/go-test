// controllers/product_controller.go
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

type ProductController struct {
	Collection *mongo.Collection
}

func NewProductController(db *mongo.Database) *ProductController {
	return &ProductController{
		Collection: db.Collection("products"),
	}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingProduct := models.Product{}
	err := pc.Collection.FindOne(context.Background(), bson.M{"code": product.Code}).Decode(&existingProduct)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Code already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = pc.Collection.InsertOne(context.Background(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}


func (pc *ProductController) GetProducts(c *gin.Context) {
	var products []models.Product
	cursor, err := pc.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product models.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

func (pc *ProductController) GetProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var product models.Product
	err = pc.Collection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = pc.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": productID},
		bson.M{"$set": updatedProduct},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var existingProduct models.Product
	err = pc.Collection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&existingProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	_, err = pc.Collection.DeleteOne(context.Background(), bson.M{"_id": productID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}


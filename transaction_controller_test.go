package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mifaabiyyu/go-test.git/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode) // Set Gin to test mode

	r := gin.Default()
	r.GET("/transactions")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Add more assertions as needed
}

func TestCreateTransaction(t *testing.T) {
	// Define the JSON data for the new transaction
	newEntry := struct {
		Transaction models.Transaction          `json:"transaction"`
		Details     []models.TransactionDetail  `json:"details"`
		Payments    []models.TransactionPayment `json:"payments"`
	}{
		Transaction: models.Transaction{
			CustomerID: primitive.NewObjectID(),
		},
		Details: []models.TransactionDetail{
			{
				ProductID: primitive.NewObjectID(),
				Price:     12000,
				Quantity:  2,
				Subtotal:  24000,
			},
			{
				ProductID: primitive.NewObjectID(),
				Price:     15000,
				Quantity:  2,
				Subtotal:  30000,
			},
		},
		Payments: []models.TransactionPayment{
			{
				PaymentMethodID: primitive.NewObjectID(),
				Status:          0,
				PaidAmount:      0,
			},
		},
	}

	r := gin.Default()

	r.POST("/transaction")
	jsonData, err := json.Marshal(newEntry)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transaction", bytes.NewBuffer(jsonData))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Add more assertions as needed
}

func TestGetTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode) // Set Gin to test mode

	idData := "64dc291883f42d242a656d61"
	r := gin.Default()
	r.GET("/transaction/" + idData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transaction/"+idData, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Add more assertions as needed
}

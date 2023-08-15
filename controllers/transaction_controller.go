// controllers/transaction_controller.go
package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mifaabiyyu/go-test.git/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionController struct {
	Collection             *mongo.Collection
	TransactionPaymentCtrl *TransactionPaymentController
}

type CreateTransactionResponse struct {
	Transaction models.Transaction `json:"transaction"`
}

func NewTransactionController(db *mongo.Database, tpc *TransactionPaymentController) *TransactionController {
	return &TransactionController{
		Collection:             db.Collection("transactions"),
		TransactionPaymentCtrl: tpc,
	}
}

func NewTransactionDetailController(db *mongo.Database) *TransactionController {
	return &TransactionController{
		Collection: db.Collection("transaction_details"),
	}
}

func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var transactionData struct {
		Transaction models.Transaction          `json:"transaction" `
		Details     []models.TransactionDetail  `json:"details" binding:"required"`
		Payments    []models.TransactionPayment `json:"payments"`
	}

	if err := c.ShouldBindJSON(&transactionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var totalAmount float64
	var totalQty float64
	for _, detail := range transactionData.Details {
		totalAmount += detail.Subtotal
		totalQty += detail.Quantity
	}

	transactionData.Transaction.TransactionDate = time.Now()
	transactionData.Transaction.TotalAmount = totalAmount
	transactionData.Transaction.TotalQty = totalQty

	ctx := context.TODO() // Create a context

	session, err := tc.Collection.Database().Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer session.EndSession(ctx)

	// Start a transaction
	err = session.StartTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert transaction
	transactionData.Transaction.ID = primitive.NewObjectID()
	_, err = tc.Collection.InsertOne(ctx, transactionData.Transaction)
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert transaction details
	for i, detail := range transactionData.Details {
		detail.ID = primitive.NewObjectID()
		detail.TransactionID = transactionData.Transaction.ID
		transactionData.Details[i] = detail
		_, err = tc.Collection.Database().Collection("transaction_details").InsertOne(ctx, detail)
		if err != nil {
			session.AbortTransaction(ctx)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Assign the transaction ID to the payment
	for i, payment := range transactionData.Payments {
		payment.ID = primitive.NewObjectID()
		payment.TransactionID = transactionData.Transaction.ID
		transactionData.Payments[i] = payment
		if err := tc.TransactionPaymentCtrl.CreateTransactionPayment(&payment); err != nil {
			session.AbortTransaction(ctx)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	// Commit the transaction
	err = session.CommitTransaction(ctx)
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := CreateTransactionResponse{
		Transaction: transactionData.Transaction,
	}

	response.Transaction.Details = transactionData.Details
	response.Transaction.Payments = transactionData.Payments

	c.JSON(http.StatusCreated, response)
}

func (tc *TransactionController) GetTransactions(c *gin.Context) {
	var transactions []models.Transaction

	orderParam := c.DefaultQuery("order_by", "transaction_date") // Default order by transaction_date
	orderDirection := 1
	if c.DefaultQuery("order_direction", "asc") == "desc" {
		orderDirection = -1
	}

	sortBy := bson.D{{Key: orderParam, Value: orderDirection}}

	cursor, err := tc.Collection.Find(context.Background(), bson.M{}, options.Find().SetSort(sortBy))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var transaction models.Transaction
		cursor.Decode(&transaction)

		// Fetch details associated with the transaction
		detailsCursor, err := tc.Collection.Database().Collection("transaction_details").Find(context.Background(), bson.M{"transaction_id": transaction.ID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer detailsCursor.Close(context.Background())

		var details []models.TransactionDetail
		for detailsCursor.Next(context.Background()) {
			var detail models.TransactionDetail
			detailsCursor.Decode(&detail)
			details = append(details, detail)
		}
		transaction.Details = details

		// Fetch payments associated with the transaction
		paymentsCursor, err := tc.Collection.Database().Collection("transaction_payments").Find(context.Background(), bson.M{"transaction_id": transaction.ID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer paymentsCursor.Close(context.Background())

		var payments []models.TransactionPayment
		for paymentsCursor.Next(context.Background()) {
			var payment models.TransactionPayment
			paymentsCursor.Decode(&payment)
			payments = append(payments, payment)
		}
		transaction.Payments = payments

		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}

func (tc *TransactionController) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	transactionID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var transaction models.Transaction
	err = tc.Collection.FindOne(context.Background(), bson.M{"_id": transactionID}).Decode(&transaction)

	detailsCursor, err := tc.Collection.Database().Collection("transaction_details").Find(context.Background(), bson.M{"transaction_id": transaction.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer detailsCursor.Close(context.Background())

	var details []models.TransactionDetail
	for detailsCursor.Next(context.Background()) {
		var detail models.TransactionDetail
		detailsCursor.Decode(&detail)
		details = append(details, detail)
	}
	transaction.Details = details

	// Fetch payments associated with the transaction
	paymentsCursor, err := tc.Collection.Database().Collection("transaction_payments").Find(context.Background(), bson.M{"transaction_id": transaction.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer paymentsCursor.Close(context.Background())

	var payments []models.TransactionPayment
	for paymentsCursor.Next(context.Background()) {
		var payment models.TransactionPayment
		paymentsCursor.Decode(&payment)
		payments = append(payments, payment)
	}
	transaction.Payments = payments

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (tc *TransactionController) UpdateTransaction(c *gin.Context) {
	transactionID := c.Param("id") // Get the transaction ID from the URL parameter

	var updatedData struct {
		Transaction models.Transaction          `json:"transaction"`
		Details     []models.TransactionDetail  `json:"details"`
		Payments    []models.TransactionPayment `json:"payments"`
	}

	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var totalAmount float64
	var totalQty float64
	for _, detail := range updatedData.Details {
		totalAmount += detail.Subtotal
		totalQty += detail.Quantity
	}

	updatedData.Transaction.TransactionDate = time.Now()
	updatedData.Transaction.TotalAmount = totalAmount
	updatedData.Transaction.TotalQty = totalQty

	ctx := context.TODO() // Create a context

	session, err := tc.Collection.Database().Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer session.EndSession(ctx)

	// Start a transaction
	err = session.StartTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update transaction
	updatedData.Transaction.ID, err = primitive.ObjectIDFromHex(transactionID)
	_, err = tc.Collection.ReplaceOne(ctx, bson.M{"_id": updatedData.Transaction.ID}, updatedData.Transaction)
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete existing transaction details
	_, err = tc.Collection.Database().Collection("transaction_details").DeleteMany(ctx, bson.M{"transaction_id": updatedData.Transaction.ID})
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert new transaction details
	for i, detail := range updatedData.Details {
		detail.ID = primitive.NewObjectID()
		detail.TransactionID = updatedData.Transaction.ID
		updatedData.Details[i] = detail
		_, err = tc.Collection.Database().Collection("transaction_details").InsertOne(ctx, detail)
		if err != nil {
			session.AbortTransaction(ctx)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit the transaction
	err = session.CommitTransaction(ctx)
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}

func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	transactionID := c.Param("id") // Get the transaction ID from the URL parameter

	ctx := context.TODO() // Create a context

	session, err := tc.Collection.Database().Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer session.EndSession(ctx)

	// Start a transaction
	err = session.StartTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the transaction
	_, err = tc.Collection.DeleteOne(ctx, bson.M{"_id": transactionID})
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete associated transaction details
	_, err = tc.Collection.Database().Collection("transaction_details").DeleteMany(ctx, bson.M{"transaction_id": transactionID})
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete associated transaction payments
	_, err = tc.Collection.Database().Collection("transaction_payments").DeleteMany(ctx, bson.M{"transaction_id": transactionID})
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	err = session.CommitTransaction(ctx)
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction and associated details/payments deleted"})
}

func (tc *TransactionController) UpdateTransactionDetail(detail models.TransactionDetail) error {
	ctx := context.TODO()

	// Define the filter to identify the specific detail you want to update
	filter := bson.M{"_id": detail.ID}

	// Define the update data
	update := bson.M{
		"$set": bson.M{
			"quantity": detail.Quantity,
			"subtotal": detail.Subtotal,
			"price":    detail.Price,
		},
	}

	// Perform the update using the updateOne method
	_, err := tc.Collection.Database().Collection("transaction_details").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TransactionController) UpdateTransactionPayment(payment models.TransactionPayment) error {
	ctx := context.TODO()

	// Define the filter to identify the specific payment you want to update
	filter := bson.M{"_id": payment.ID}

	// Define the update data
	update := bson.M{
		"$set": bson.M{
			"status":       payment.Status,
			"paid_amount":  payment.PaidAmount,
			"payment_date": payment.PaymentDate,
		},
	}

	// Perform the update using the updateOne method
	_, err := tc.Collection.Database().Collection("transaction_payments").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mifaabiyyu/go-test.git/controllers"
)

func main() {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb+srv://mifaabiyyu:indun1234@cluster0.geg1dqt.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	db := client.Database("mydatabase")

	// Set up Gin router
	router := gin.Default()

	// Initialize controller
	customerController := controllers.NewCustomerController(db)
	customerAddressController := controllers.NewCustomerAddressController(db)
	productController := controllers.NewProductController(db)
	paymentMethodController := controllers.NewPaymentMethodController(db)
	transactionController := controllers.NewTransactionController(db, controllers.NewTransactionPaymentController(db))

	// Define routes
	router.POST("/customers", customerController.CreateCustomer)
	router.GET("/customers", customerController.GetCustomers)
	router.GET("/customers/:id", customerController.GetCustomer)
	router.PUT("/customers/:id", customerController.UpdateCustomer)
	router.DELETE("/customers/:id", customerController.DeleteCustomer)

	router.POST("/customer-addresses", customerAddressController.CreateCustomerAddress)
	router.GET("/customer-addresses", customerAddressController.GetCustomerAddresses)
	router.GET("/customer-addresses/:id", customerAddressController.GetCustomerAddress)
	router.PUT("/customer-addresses/:id", customerAddressController.UpdateCustomerAddress)
	router.DELETE("/customer-addresses/:id", customerAddressController.DeleteCustomerAddress)

	router.GET("/products", productController.GetProducts)
	router.POST("/product", productController.CreateProduct)
	router.GET("/product/:id", productController.GetProduct)
	router.PUT("/product/:id", productController.UpdateProduct)
	router.DELETE("/product/:id", productController.DeleteProduct)

	router.GET("/payment-methods", paymentMethodController.GetPaymentMethods)
	router.POST("/payment-method", paymentMethodController.CreatePaymentMethod)
	router.GET("/payment-method/:id", paymentMethodController.GetPaymentMethod)
	router.PUT("/payment-method/:id", paymentMethodController.UpdatePaymentMethod)
	router.DELETE("/payment-method/:id", paymentMethodController.DeletePaymentMethod)

	router.GET("/transactions", transactionController.GetTransactions)
	router.GET("/transaction/:id", transactionController.GetTransaction)
	router.POST("/transaction", transactionController.CreateTransaction)
	router.PUT("/transaction/:id", transactionController.UpdateTransaction)
	router.DELETE("/transaction/:id", transactionController.DeleteTransaction)
	// Start the server
	router.Run(":8080")
}

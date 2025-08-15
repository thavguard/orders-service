package main

import (
	"context"
	"log"
	"orders/src/db"
	"orders/src/db/queries"
	usecases "orders/src/use-cases"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx := context.Background()

	// Init DB

	connString := os.Getenv("DATABASE_URL")

	db, err := db.InitDbConnection(ctx, connString)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	dbService := &queries.DBService{DB: db}

	// deliveryDto := &models.Delivery{
	// 	Name:    "name",
	// 	Phone:   "892358585",
	// 	Zip:     "12321",
	// 	City:    "sdasd",
	// 	Address: "dsadasd",
	// 	Region:  "region",
	// 	Email:   "test@test.com",
	// }

	// result, createError := dbService.GetDeliveryById(ctx, 1)

	usecases.ListenOrders(ctx, dbService)

	defer db.Close()

}

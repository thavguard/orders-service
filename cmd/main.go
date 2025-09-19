package main

import (
	"context"
	"fmt"
	"log"
	"orders/src/broker/consumers"
	"orders/src/db"
	"orders/src/db/repositories"
	httpserver "orders/src/http-server"
	"orders/src/mycache"
	"orders/src/service"
	customvalidator "orders/src/utils/custom-validator"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var Validate *validator.Validate

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

func main() {
	// Создаем контекст для безопастного завершения работы
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	Validate, err := customvalidator.NewValidator()

	if err != nil {
		log.Printf("Error in create validator: %v\n", err)
	}

	// Инициализация БД
	db, err := db.NewDBConnection(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Инициализация репозиториев
	orderRepo := repositories.NewOrderRepo(db.Pool)
	itemRepo := repositories.NewItemRepo(db.Pool)
	paymentRepo := repositories.NewPaymentRepo(db.Pool)
	deliveryRepo := repositories.NewDeliveryRepo(db.Pool)

	// Инициализация redis
	redis := mycache.NewRedis(time.Minute)

	// Инициализация сервисов
	ordersService := service.NewOrderService(orderRepo, redis, Validate)
	itemService := service.NewItemService(itemRepo, redis, Validate)
	paymentService := service.NewPaymentService(paymentRepo, redis, Validate)
	deliveryService := service.NewDeliveryService(deliveryRepo, redis, Validate)

	// Создание web-server
	srv := httpserver.NewServer(ctx, ordersService)

	// Подписка на топик
	listener := consumers.NewOrderConsumer(ordersService, deliveryService, itemService, paymentService)
	listener.Run(ctx)

	// Обработка закрытия  приложения
	<-ctx.Done()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error in srv.Shutdown: %v\n", err)
	}

	if err1, err2 := listener.Close(); err != nil || err2 != nil {
		if err1 != nil {
			log.Printf("Error in reader.Close: %v\n", err1)
		}

		if err2 != nil {
			log.Printf("Error in writer.Close: %v\n", err2)
		}
	}

	if err := redis.Close(); err != nil {
		log.Printf("Error in redis.Close: %v\n", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("Error in db.Close: %v\n", err)
	}

}

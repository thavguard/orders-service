package main

import (
	"context"
	"log"
	"orders/src/broker/consumers"
	"orders/src/db"
	"orders/src/db/repositories"
	httpserver "orders/src/http-server"
	"orders/src/metrics"
	"orders/src/mycache"
	"orders/src/service"
	"orders/src/tracer"
	customvalidator "orders/src/utils/custom-validator"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
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

	// Tracing
	tp, err := tracer.InitTracer(os.Getenv("JAEGER_URL"), "Orders Service")

	if err != nil {
		log.Printf("init tracer: %v\n", err)
	}

	// Метрики
	reg := prometheus.DefaultRegisterer
	gt := prometheus.DefaultGatherer

	met := metrics.New(reg, gt)

	Validate, err := customvalidator.NewValidator()

	if err != nil {
		log.Printf("Error in create validator: %v\n", err)
	}

	// Инициализация БД
	db, err := db.NewDBConnection(ctx, tp, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Инициализация репозиториев
	orderRepo := repositories.NewOrderRepo(db.Pool, met)
	itemRepo := repositories.NewItemRepo(db.Pool, met)
	paymentRepo := repositories.NewPaymentRepo(db.Pool, met)
	deliveryRepo := repositories.NewDeliveryRepo(db.Pool, met)

	// Инициализация redis
	redis := mycache.NewRedis(tp, reg, met, time.Minute)

	// Инициализация сервисов
	ordersService := service.NewOrderService(orderRepo, redis, Validate)
	itemService := service.NewItemService(itemRepo, redis, Validate)
	paymentService := service.NewPaymentService(paymentRepo, redis, Validate)
	deliveryService := service.NewDeliveryService(deliveryRepo, redis, Validate)

	// Создание web-server
	srv := httpserver.NewServer(ctx, met, ordersService)

	// Подписка на топик
	listener := consumers.NewOrderConsumer(met, tp, ordersService, deliveryService, itemService, paymentService)
	listener.Run(ctx)

	// Обработка закрытия  приложения
	<-ctx.Done()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error in srv.Shutdown: %v\n", err)
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

	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Error in tp.Shutdown: %v\n", err)
	}

}

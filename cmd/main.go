package main

import (
	"context"
	"log"
	"orders/src/db"
	"orders/src/db/queries"
	httpserver "orders/src/http-server"
	usecases "orders/src/use-cases"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

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

	var wg sync.WaitGroup

	// Инициализация БД

	db, err := db.InitDbConnection(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	dbService := &queries.DBService{DB: db}

	// Подписка на топик
	reader := usecases.ListenOrders(ctx, dbService)

	httpPort := ":" + os.Getenv("HTTP_PORT")

	// Создание web-server
	srv := httpserver.NewServer(ctx, &wg, httpPort, dbService)

	// Обработка закрытия  приложения
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	reader.Close()

	wg.Wait()

}

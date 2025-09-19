package httpserver

import (
	"context"
	"log"
	"net/http"
	orderroute "orders/src/http-server/order-route"
	"orders/src/service"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(ctx context.Context, orderService *service.OrderService) *http.Server {
	httpPort := ":" + os.Getenv("HTTP_PORT")

	router := gin.Default()
	router.Use(cors.Default()) // All origins allowed by default

	orderroute.AddOrderRoutes(ctx, router, orderService)

	srv := &http.Server{
		Addr:              httpPort,
		Handler:           router.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http: %v", err)
		}

	}()

	return srv

}

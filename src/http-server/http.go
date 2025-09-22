package httpserver

import (
	"context"
	"log"
	"net/http"
	orderroute "orders/src/http-server/order-route"
	"orders/src/metrics"
	"orders/src/service"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewServer(ctx context.Context, met *metrics.Metrics, orderService *service.OrderService) *http.Server {
	httpPort := ":" + os.Getenv("HTTP_PORT")

	router := gin.Default()
	router.Use(otelgin.Middleware("http-service"))
	router.Use(cors.Default()) // All origins allowed by default

	router.Use(GinMetricsMiddleware(met))

	router.GET("/metrics", gin.WrapH(met.Handler()))

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

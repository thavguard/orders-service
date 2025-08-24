package httpserver

import (
	"context"
	"log"
	"net/http"
	"orders/src/service"
	"sync"

	"github.com/gin-gonic/gin"
)

func NewServer(ctx context.Context, wg *sync.WaitGroup, port string, service *service.Service) *http.Server {
	router := gin.Default()

	AddRoutes(ctx, service, router)

	srv := &http.Server{
		Addr:    port,
		Handler: router.Handler(),
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("http: %v", err)
		}

	}()

	return srv

}

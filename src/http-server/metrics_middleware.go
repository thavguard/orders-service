package httpserver

import (
	"fmt"
	"time"

	"orders/src/metrics"

	"github.com/gin-gonic/gin"
)

func GinMetricsMiddleware(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		m.HTTPInflight.Inc()
		defer m.HTTPInflight.Dec()

		c.Next() // выполняем обработчик

		status := c.Writer.Status()
		duration := time.Since(start).Seconds()
		handlerName := c.FullPath()
		if handlerName == "" {
			handlerName = "unknown"
		}

		m.HTTPRequestCount.WithLabelValues(handlerName, c.Request.Method,
			fmt.Sprintf("%d", status), "orders_service").Inc()
		m.HTTPRequestDuration.WithLabelValues(handlerName, c.Request.Method, "orders_service").Observe(duration)
	}
}

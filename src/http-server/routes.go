package httpserver

import (
	"context"
	"orders/src/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddRoutes(ctx context.Context, service *service.Service, router *gin.Engine) {
	router.GET("/order/:orderId", func(c *gin.Context) {
		orderId, err := strconv.Atoi(c.Param("orderId"))

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "orderId must be an integer string",
			})
			return
		}

		order, err := service.GetOrderById(ctx, orderId)

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"message": err,
			})
			return

		}

		c.JSON(200, gin.H{
			"order": order,
		})

	})
}

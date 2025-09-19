package orderroute

import (
	"context"
	"fmt"
	"orders/src/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddOrderRoutes(ctx context.Context, router *gin.Engine, orderService *service.OrderService) {

	router.GET("/order/:orderID", func(c *gin.Context) {
		orderID, err := strconv.Atoi(c.Param("orderID"))

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "orderID must be an integer string",
			})
			return
		}
		order, err := orderService.GetOrderByID(ctx, orderID)

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"message": fmt.Sprintf("Error: %v\n", err),
			})
			return

		}
		c.JSON(200, gin.H{
			"order": order,
		})
	})

}

package service

import (
	"context"
	"fmt"
	"orders/src/internal/broker"
	"strconv"
	"time"

	"github.com/go-redis/cache/v9"
)

func (service *Service) GetOrderById(ctx context.Context, orderId int) (*broker.OrderMessage, error) {
	redisKey := "order_" + strconv.Itoa(orderId)

	var order *broker.OrderMessage

	err := service.MyCache.Cache.Get(ctx, redisKey, &order)

	if err != nil {
		fmt.Printf("ERROR IN REDIS: %v\n", err)

		order, err = service.Repo.GetOrderById(ctx, orderId)

		if err != nil {
			fmt.Printf("ERROR IN DB: %v\n", err)
		} else {
			service.MyCache.Cache.Set(&cache.Item{
				Ctx:   ctx,
				Key:   redisKey,
				Value: &order,
				TTL:   time.Hour,
			})
		}
	}

	return order, err
}

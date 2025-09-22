package service

import (
	"context"
	"log"
	"orders/src/broker"
	"orders/src/db/models"
	"orders/src/db/repositories"
	"orders/src/mycache"
	"strconv"

	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/singleflight"
)

type OrderService struct {
	myCache   *mycache.RedisService
	orderRepo repositories.OrderRepository
	valid     *validator.Validate
	g         singleflight.Group
}

func NewOrderService(
	orderRepo repositories.OrderRepository, myCache *mycache.RedisService, valid *validator.Validate) *OrderService {
	return &OrderService{myCache: myCache, orderRepo: orderRepo, valid: valid}
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID int) (*broker.OrderMessage, error) {
	redisKey := "order_" + strconv.Itoa(orderID)

	var order *broker.OrderMessage

	err := s.myCache.Get(ctx, redisKey, &order)

	if err != nil {
		log.Printf("ERROR IN REDIS: %v\n", err)

		v, err, _ := s.g.Do(redisKey, func() (interface{}, error) {
			return s.orderRepo.GetOrderByID(ctx, orderID)

		})

		if err != nil {
			log.Printf("ERROR IN DB: %v\n", err)
			return order, err
		}

		order = v.(*broker.OrderMessage)

		if err = s.myCache.Set(ctx, redisKey, &order); err != nil {
			log.Printf("Error in Cache Set %v\n", err)
		}

		s.g.Forget(redisKey)
	}

	return order, err
}

func (s *OrderService) CreateOrder(ctx context.Context, orderDto models.Order) (models.Order, error) {
	if err := s.valid.StructCtx(ctx, orderDto); err != nil {
		return models.Order{}, err
	}

	order, err := s.orderRepo.CreateOrder(ctx, &orderDto)

	if err != nil {
		log.Printf("Error in CreateOrder: %v\n", err)
		return models.Order{}, err
	}

	return order, nil

}

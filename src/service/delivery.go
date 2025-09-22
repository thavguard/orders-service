package service

import (
	"context"
	"log"
	"orders/src/db/models"
	"orders/src/db/repositories"
	"orders/src/mycache"
	"strconv"

	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/singleflight"
)

type DeliveryService struct {
	myCache      *mycache.RedisService
	deliveryRepo repositories.DeliveryRepository
	valid        *validator.Validate
	g            singleflight.Group
}

func NewDeliveryService(repo repositories.DeliveryRepository, cache *mycache.RedisService, valid *validator.Validate) *DeliveryService {
	return &DeliveryService{deliveryRepo: repo, myCache: cache, valid: valid}
}

func (s *DeliveryService) CreateDelivery(ctx context.Context, deliveryDto *models.Delivery) (models.Delivery, error) {
	if err := s.valid.StructCtx(ctx, deliveryDto); err != nil {
		log.Printf("ERROR IN VALIDATE: %v\n", err)

		return models.Delivery{}, err
	}

	delivery, err := s.deliveryRepo.CreateDelivery(ctx, deliveryDto)

	if err != nil {
		log.Printf("ERROR IN CreateDelivery: %v\n", err)

		return models.Delivery{}, err
	}

	return delivery, nil

}

func (s *DeliveryService) GetDeliveryByID(ctx context.Context, deliveryID int) (models.Delivery, error) {

	var delivery models.Delivery

	redisKey := "delivery_" + strconv.Itoa(deliveryID)

	if err := s.myCache.Get(ctx, redisKey, &delivery); err != nil {

		v, err, _ := s.g.Do(redisKey, func() (interface{}, error) {
			return s.deliveryRepo.GetDeliveryByID(ctx, deliveryID)
		})

		if err != nil {
			return models.Delivery{}, err
		}

		delivery = v.(models.Delivery)

		if err = s.myCache.Set(ctx, redisKey, delivery); err != nil {
			log.Printf("ERROR IN SET CACHE: %v\n", err)
		}

		s.g.Forget(redisKey)

	}

	return delivery, nil
}

func (s *DeliveryService) GetDeliveryByOrderID(ctx context.Context, orderID int) (models.Delivery, error) {

	var delivery models.Delivery

	redisKey := "devivery_order_id" + strconv.Itoa(orderID)

	if err := s.myCache.Get(ctx, redisKey, &delivery); err != nil {

		v, err, _ := s.g.Do(redisKey, func() (interface{}, error) {
			return s.deliveryRepo.GetDeliveryByOrderID(ctx, orderID)
		})

		if err != nil {
			return models.Delivery{}, err
		}
		delivery := v.(models.Delivery)

		if err = s.myCache.Set(ctx, redisKey, delivery); err != nil {
			log.Printf("ERROR IN SET CACHE: %v\n", err)
		}

		s.g.Forget(redisKey)

	}

	return delivery, nil
}

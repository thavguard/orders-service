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

type PaymentService struct {
	myCache     *mycache.RedisService
	paymentRepo repositories.PaymentRepository
	valid       *validator.Validate
	g           singleflight.Group
}

func NewPaymentService(repo repositories.PaymentRepository, cache *mycache.RedisService, valid *validator.Validate) *PaymentService {
	return &PaymentService{paymentRepo: repo, myCache: cache, valid: valid}
}

func (s *PaymentService) CreatePayment(ctx context.Context, paymentDto *models.Payment) (models.Payment, error) {
	if err := s.valid.StructCtx(ctx, paymentDto); err != nil {
		log.Printf("ERROR IN VALIDATE: %v\n", err)

		return models.Payment{}, err
	}

	payment, err := s.paymentRepo.CreatePayment(ctx, paymentDto)

	if err != nil {
		log.Printf("ERROR IN CreatePayment: %v\n", err)

		return models.Payment{}, err
	}

	return payment, nil

}

func (s *PaymentService) GetPaymentByID(ctx context.Context, paymentID int) (models.Payment, error) {

	var payment models.Payment

	redisKey := "payment_" + strconv.Itoa(paymentID)

	if err := s.myCache.Get(ctx, redisKey, &payment); err != nil {

		v, err, _ := s.g.Do(redisKey, func() (interface{}, error) {
			return s.paymentRepo.GetPaymentByID(ctx, paymentID)
		})

		if err != nil {
			return models.Payment{}, err
		}

		payment = v.(models.Payment)

		if err = s.myCache.Set(ctx, redisKey, payment); err != nil {
			log.Printf("ERROR IN SET CACHE: %v\n", err)
		}

		s.g.Forget(redisKey)

	}

	return payment, nil
}

func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID int) (models.Payment, error) {

	var payment models.Payment

	redisKey := "payment_order_id" + strconv.Itoa(orderID)

	if err := s.myCache.Get(ctx, redisKey, &payment); err != nil {

		v, err, _ := s.g.Do(redisKey, func() (interface{}, error) {
			return s.paymentRepo.GetPaymentByOrderID(ctx, orderID)
		})

		if err != nil {
			return models.Payment{}, err
		}

		payment = v.(models.Payment)

		if err = s.myCache.Set(ctx, redisKey, payment); err != nil {
			log.Printf("ERROR IN SET CACHE: %v\n", err)
		}

		s.g.Forget(redisKey)

	}

	return payment, nil
}

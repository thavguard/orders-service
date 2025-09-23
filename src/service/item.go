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

type ItemService interface {
	CreateItem(ctx context.Context, itemDto *models.Item) (models.Item, error)
	GetItemByID(ctx context.Context, itemID int) (models.Item, error)
	GetItemsByOrderID(ctx context.Context, orderID int) ([]models.Item, error)
}

type itemService struct {
	myCache  mycache.CacheService
	itemRepo repositories.ItemRepository
	valid    *validator.Validate
	g        singleflight.Group
}

func NewItemService(repo repositories.ItemRepository, cache mycache.CacheService, valid *validator.Validate) ItemService {

	return &itemService{itemRepo: repo, myCache: cache, valid: valid}
}

func (s *itemService) CreateItem(ctx context.Context, itemDto *models.Item) (models.Item, error) {
	if err := s.valid.StructCtx(ctx, itemDto); err != nil {
		log.Printf("ERROR IN VALIDATE: %v\n", err)

		return models.Item{}, err
	}

	item, err := s.itemRepo.CreateItem(ctx, itemDto)

	if err != nil {
		log.Printf("ERROR IN CreateItem: %v\n", err)

		return models.Item{}, err
	}

	return item, nil

}

func (s *itemService) GetItemByID(ctx context.Context, itemID int) (models.Item, error) {

	var item models.Item

	redisKey := "item_" + strconv.Itoa(itemID)

	if err := s.myCache.Get(ctx, redisKey, &item); err != nil {

		v, err, _ := s.g.Do(redisKey, func() (interface{}, error) {
			return s.itemRepo.GetItemByID(ctx, itemID)
		})

		if err != nil {
			return models.Item{}, err
		}

		item := v.(models.Item)

		if err = s.myCache.Set(ctx, redisKey, item); err != nil {
			log.Printf("ERROR IN PAYMENT CACHE SET: %v\n", err)
		}

		s.g.Forget(redisKey)

	}

	return item, nil
}

func (s *itemService) GetItemsByOrderID(ctx context.Context, orderID int) ([]models.Item, error) {

	var items []models.Item

	redisKey := "items_order_id" + strconv.Itoa(orderID)

	if err := s.myCache.Get(ctx, redisKey, &items); err != nil {

		v, err, _ := s.g.Do(redisKey, func() (interface{}, error) {
			return s.itemRepo.GetItemsByOrderID(ctx, orderID)
		})

		if err != nil {
			return []models.Item{}, err
		}

		items = v.([]models.Item)

		if err = s.myCache.Set(ctx, redisKey, items); err != nil {
			log.Printf("ERROR IN PAYMENT CACHE SET: %v\n", err)
		}

		s.g.Forget(redisKey)

	}

	return items, nil
}

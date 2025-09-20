package repositories

import (
	"context"
	"fmt"
	"orders/src/db"
	"orders/src/db/models"
	"orders/src/metrics"
	"orders/src/myretry"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-retry"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, itemDto *models.Item) (models.Item, error)
	GetItemByID(ctx context.Context, itemID int) (models.Item, error)
	GetItemsByOrderID(ctx context.Context, orderID int) ([]models.Item, error)
}

type itemRepo struct {
	pool    *sqlx.DB
	b       func() retry.Backoff
	metrics *metrics.Metrics
}

func NewItemRepo(pool *sqlx.DB, metrics *metrics.Metrics) ItemRepository {
	b := myretry.NewBackofFactory()
	return &itemRepo{pool: pool, b: b, metrics: metrics}
}

func (repo *itemRepo) CreateItem(ctx context.Context, itemDto *models.Item) (models.Item, error) {
	var item models.Item
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {
		item, err = repo.createItem(ctx, itemDto)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return item, err
}

func (repo *itemRepo) createItem(ctx context.Context, itemDto *models.Item) (models.Item, error) {
	start := time.Now()

	var item models.Item

	query := `
     INSERT INTO item (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id)
VALUES (:chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status, :order_id)
RETURNING id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id;
    `

	rows, err := repo.pool.NamedQueryContext(ctx, query, &itemDto)

	lat := time.Since(start).Seconds()
	repo.metrics.DBQueryDuration.WithLabelValues("create_item", "item_service").Observe(lat)

	if err != nil {
		repo.metrics.DBQueryErrors.WithLabelValues("create_item", "item_service").Inc()

		return models.Item{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.StructScan(&item); err != nil {
			fmt.Printf("Error while parsing rows %v", err)
			return models.Item{}, err
		}
	}

	return item, nil
}

func (repo *itemRepo) GetItemByID(ctx context.Context, itemID int) (models.Item, error) {
	var item models.Item
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		item, err = repo.getItemByID(ctx, itemID)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return item, err
}

func (repo *itemRepo) getItemByID(ctx context.Context, itemID int) (models.Item, error) {

	start := time.Now()

	var item models.Item

	query := `select *
			from item
			where id = $1;`

	err := repo.pool.GetContext(ctx, &item, query, itemID)

	lat := time.Since(start).Seconds()
	repo.metrics.DBQueryDuration.WithLabelValues("get_item_by_id", "item_service").Observe(lat)

	if err != nil {
		repo.metrics.DBQueryErrors.WithLabelValues("get_item_by_id", "item_service").Inc()

		fmt.Printf("Error in GetItemByID: %v\n", err)
		return item, err
	}

	return item, nil
}

func (repo *itemRepo) GetItemsByOrderID(ctx context.Context, orderID int) ([]models.Item, error) {
	var items []models.Item
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {
		items, err = repo.getItemsByOrderID(ctx, orderID)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err

	})

	return items, err
}

func (repo *itemRepo) getItemsByOrderID(ctx context.Context, orderID int) ([]models.Item, error) {

	start := time.Now()

	var items []models.Item

	query := `select *
			from item
			where order_id = $1 order by id;`

	err := repo.pool.SelectContext(ctx, &items, query, orderID)

	lat := time.Since(start).Seconds()
	repo.metrics.DBQueryDuration.WithLabelValues("get_item_by_order_id", "item_service").Observe(lat)

	if err != nil {
		repo.metrics.DBQueryErrors.WithLabelValues("get_item_by_order_id", "item_service").Inc()

		fmt.Printf("Error in GetItemsByOrderID: %v\n", err)
		return items, err
	}

	return items, nil
}

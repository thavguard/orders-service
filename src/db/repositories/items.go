package repositories

import (
	"context"
	"fmt"
	"orders/src/db"
	"orders/src/db/models"
	"orders/src/myretry"

	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-retry"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, itemDto *models.Item) (models.Item, error)
	GetItemByID(ctx context.Context, itemID int) (models.Item, error)
	GetItemsByOrderID(ctx context.Context, orderID int) ([]models.Item, error)
}

type itemRepo struct {
	pool *sqlx.DB
	b    func() retry.Backoff
}

func NewItemRepo(pool *sqlx.DB) ItemRepository {
	b := myretry.NewBackofFactory()
	return &itemRepo{pool: pool, b: b}
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
	var item models.Item

	query := `
     INSERT INTO item (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id)
VALUES (:chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status, :order_id)
RETURNING id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id;
    `

	rows, err := repo.pool.NamedQueryContext(ctx, query, &itemDto)

	if err != nil {
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
	var item models.Item

	query := `select *
			from item
			where id = $1;`

	err := repo.pool.GetContext(ctx, &item, query, itemID)

	if err != nil {
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
	var items []models.Item

	query := `select *
			from item
			where order_id = $1 order by id;`

	err := repo.pool.SelectContext(ctx, &items, query, orderID)

	if err != nil {
		fmt.Printf("Error in GetItemsByOrderID: %v\n", err)
		return items, err
	}

	return items, nil
}

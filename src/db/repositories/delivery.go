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

type DeliveryRepository interface {
	CreateDelivery(ctx context.Context, deliveryDto *models.Delivery) (models.Delivery, error)
	GetDeliveryByID(ctx context.Context, deliveryID int) (models.Delivery, error)
	GetDeliveryByOrderID(ctx context.Context, orderID int) (models.Delivery, error)
}

type deliveryRepo struct {
	pool *sqlx.DB
	b    func() retry.Backoff
}

func NewDeliveryRepo(pool *sqlx.DB) DeliveryRepository {
	b := myretry.NewBackofFactory()
	return &deliveryRepo{pool: pool, b: b}
}

func (repo *deliveryRepo) CreateDelivery(ctx context.Context, deliveryDto *models.Delivery) (models.Delivery, error) {
	var delivery models.Delivery
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		delivery, err = repo.createDelivery(ctx, deliveryDto)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return delivery, err
}

func (repo *deliveryRepo) createDelivery(ctx context.Context, deliveryDto *models.Delivery) (models.Delivery, error) {
	var delivery models.Delivery

	query := `
     INSERT INTO delivery (name, phone, zip, city, address, region, email)
VALUES (:name, :phone, :zip, :city, :address, :region, :email)
RETURNING id, name, phone, zip, city, address, region, email;
    `

	rows, err := repo.pool.NamedQueryContext(ctx, query, &deliveryDto)

	if err != nil {
		return models.Delivery{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.StructScan(&delivery); err != nil {
			fmt.Printf("Error while parsing  rows %v", err)
			return models.Delivery{}, err
		}
	}

	return delivery, nil
}

func (repo *deliveryRepo) GetDeliveryByID(ctx context.Context, deliveryID int) (models.Delivery, error) {
	var delivery models.Delivery
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		delivery, err = repo.getDeliveryByID(ctx, deliveryID)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return delivery, err
}

func (repo *deliveryRepo) getDeliveryByID(ctx context.Context, deliveryID int) (models.Delivery, error) {
	var delivery models.Delivery

	query := `select *
			from delivery
			where id = $1;`

	err := repo.pool.GetContext(ctx, &delivery, query, deliveryID)

	if err != nil {
		fmt.Printf("Error in GetDeliveryByID: %v\n", err)
		return models.Delivery{}, err
	}

	return delivery, nil
}

func (repo *deliveryRepo) GetDeliveryByOrderID(ctx context.Context, orderID int) (models.Delivery, error) {
	var delivery models.Delivery
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		delivery, err = repo.getDeliveryByOrderID(ctx, orderID)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return delivery, err
}

func (repo *deliveryRepo) getDeliveryByOrderID(ctx context.Context, orderID int) (models.Delivery, error) {
	var delivery models.Delivery

	query := `select *
			from delivery
			where order_id = $1;`

	err := repo.pool.GetContext(ctx, &delivery, query, orderID)

	if err != nil {
		fmt.Printf("Error in GetDeliveryByOrderID: %v\n", err)
		return models.Delivery{}, err
	}

	return delivery, nil
}

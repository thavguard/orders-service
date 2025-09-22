package repositories

import (
	"context"
	"log"
	"orders/src/db"
	"orders/src/db/models"
	"orders/src/metrics"
	"orders/src/myretry"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-retry"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, paymentDto *models.Payment) (models.Payment, error)
	GetPaymentByID(ctx context.Context, paymentID int) (models.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID int) (models.Payment, error)
}

type paymentRepo struct {
	pool    *sqlx.DB
	b       func() retry.Backoff
	metrics *metrics.Metrics
}

func NewPaymentRepo(pool *sqlx.DB, metrics *metrics.Metrics) PaymentRepository {
	b := myretry.NewBackofFactory()
	return &paymentRepo{pool: pool, b: b, metrics: metrics}
}

func (repo *paymentRepo) CreatePayment(ctx context.Context, paymentDto *models.Payment) (models.Payment, error) {
	var payment models.Payment
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		payment, err = repo.createPayment(ctx, paymentDto)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return payment, err
}

func (repo *paymentRepo) createPayment(ctx context.Context, paymentDto *models.Payment) (models.Payment, error) {
	start := time.Now()

	var payment models.Payment

	query := `
     INSERT INTO payment (currency, delivery_cost, provider,
amount, payment_dt, bank, request_id, transaction, custom_fee, goods_total)
VALUES (:currency, :delivery_cost, :provider,
:amount, :payment_dt, :bank, :request_id, :transaction, :custom_fee, :goods_total)
RETURNING id, currency, delivery_cost, provider,
amount, payment_dt, bank, request_id, transaction, custom_fee, goods_total;
    `

	rows, err := repo.pool.NamedQueryContext(ctx, query, &paymentDto)

	lat := time.Since(start).Seconds()
	repo.metrics.DBQueryDuration.WithLabelValues("create_payment", "payment_service").Observe(lat)

	if err != nil {
		repo.metrics.DBQueryErrors.WithLabelValues("create_payment", "payment_service").Inc()

		return models.Payment{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.StructScan(&payment); err != nil {
			log.Printf("Error while parsing rows %v", err)
			return models.Payment{}, err
		}
	}

	return payment, nil
}

func (repo *paymentRepo) GetPaymentByID(ctx context.Context, paymentID int) (models.Payment, error) {
	var payment models.Payment
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		payment, err = repo.getPaymentByID(ctx, paymentID)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return payment, err
}

func (repo *paymentRepo) getPaymentByID(ctx context.Context, paymentID int) (models.Payment, error) {
	start := time.Now()

	var payment models.Payment

	query := `select *
			from payment
			where id = $1;`

	err := repo.pool.GetContext(ctx, &payment, query, paymentID)

	lat := time.Since(start).Seconds()
	repo.metrics.DBQueryDuration.WithLabelValues("get_payment_by_id", "payment_service").Observe(lat)

	if err != nil {
		repo.metrics.DBQueryErrors.WithLabelValues("get_payment_by_id", "payment_service").Inc()

		log.Printf("Error in GetPaymentByID: %v\n", err)
		return models.Payment{}, err
	}

	return payment, nil
}

func (repo *paymentRepo) GetPaymentByOrderID(ctx context.Context, orderID int) (models.Payment, error) {
	var payment models.Payment
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		payment, err = repo.getPaymentByOrderID(ctx, orderID)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return payment, err
}

func (repo *paymentRepo) getPaymentByOrderID(ctx context.Context, orderID int) (models.Payment, error) {
	start := time.Now()

	var payment models.Payment

	query := `select *
			from payment
			where order_id = $1;`

	err := repo.pool.GetContext(ctx, &payment, query, orderID)

	lat := time.Since(start).Seconds()
	repo.metrics.DBQueryDuration.WithLabelValues("get_payment_by_order_id", "payment_service").Observe(lat)

	if err != nil {
		repo.metrics.DBQueryErrors.WithLabelValues("get_payment_by_order_id", "payment_service").Inc()

		log.Printf("Error in GetPaymentByID: %v\n", err)
		return models.Payment{}, err
	}

	return payment, nil
}

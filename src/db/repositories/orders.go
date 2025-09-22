package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"orders/src/broker"
	"orders/src/db"
	"orders/src/db/models"
	"orders/src/metrics"
	"orders/src/myretry"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-retry"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, orderDto *models.Order) (models.Order, error)
	GetOrderByID(ctx context.Context, orderID int) (*broker.OrderMessage, error)
}

type orderRepo struct {
	pool    *sqlx.DB
	b       func() retry.Backoff
	metrics *metrics.Metrics
}

func NewOrderRepo(pool *sqlx.DB, metrics *metrics.Metrics) OrderRepository {
	b := myretry.NewBackofFactory()
	return &orderRepo{pool: pool, b: b, metrics: metrics}
}

func (repo *orderRepo) CreateOrder(ctx context.Context, orderDto *models.Order) (models.Order, error) {

	var order models.Order
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		order, err = repo.createOrder(ctx, orderDto)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return order, err
}

func (repo *orderRepo) createOrder(ctx context.Context, orderDto *models.Order) (models.Order, error) {

	start := time.Now()

	var order models.Order

	query := `
    INSERT INTO "order" (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created,
	oof_shard, delivery_id, payment_id)
	VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created,
	:oof_shard, :delivery_id, :payment_id)
	RETURNING id, order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created,
	oof_shard, delivery_id, payment_id;
    `

	rows, err := repo.pool.NamedQueryContext(ctx, query, orderDto)

	lat := time.Since(start).Seconds()
	repo.metrics.DBQueryDuration.WithLabelValues("create_order", "order_service").Observe(lat)

	if err != nil {
		repo.metrics.DBQueryErrors.WithLabelValues("create_order", "order_service").Inc()

		return models.Order{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.StructScan(&order); err != nil {
			log.Printf("Error while parsing rows %v", err)
			return models.Order{}, err
		}
	}

	return order, nil
}

func (repo *orderRepo) GetOrderByID(ctx context.Context, orderID int) (*broker.OrderMessage, error) {
	var order *broker.OrderMessage
	var err error

	err = retry.Do(ctx, repo.b(), func(ctx context.Context) error {

		order, err = repo.getOrderByID(ctx, orderID)

		if db.IsRetryable(err) {
			return retry.RetryableError(err)
		}

		return err
	})

	return order, err
}

func (repo *orderRepo) getOrderByID(ctx context.Context, orderID int) (*broker.OrderMessage, error) {
	start := time.Now()

	type orderRow struct {
		// Order
		OrderID           int       `db:"order_id"`
		OrderUID          string    `db:"order_order_uid"`
		TrackNumber       string    `db:"order_track_number"`
		Entry             string    `db:"order_entry"`
		Locale            string    `db:"order_locale"`
		InternalSignature string    `db:"order_internal_signature"`
		CustomerID        string    `db:"order_customer_id"`
		OrderDeliveryID   int       `db:"order_delivery_id"`
		OrderPaymentID    int       `db:"order_payment_id"`
		Shardkey          string    `db:"order_shardkey"`
		SmID              int       `db:"order_sm_id"`
		DateCreated       time.Time `db:"order_date_created"`
		OofShard          string    `db:"order_oof_shard"`

		// Payment
		PaymentID    sql.NullInt64  `db:"payment_id"`
		Transaction  sql.NullString `db:"payment_transaction"`
		RequestID    sql.NullString `db:"payment_request_id"`
		Currency     sql.NullString `db:"payment_currency"`
		Provider     sql.NullString `db:"payment_provider"`
		Amount       sql.NullInt64  `db:"payment_amount"`
		PaymentDt    sql.NullInt64  `db:"payment_payment_dt"`
		Bank         sql.NullString `db:"payment_bank"`
		DeliveryCost sql.NullInt64  `db:"payment_delivery_cost"`
		GoodsTotal   sql.NullInt64  `db:"payment_goods_total"`
		CustomFee    sql.NullInt64  `db:"payment_custom_fee"`

		// Delivery
		DeliveryID   sql.NullInt64  `db:"delivery_id"`
		DeliveryName sql.NullString `db:"delivery_name"`
		Phone        sql.NullString `db:"delivery_phone"`
		Zip          sql.NullString `db:"delivery_zip"`
		City         sql.NullString `db:"delivery_city"`
		Address      sql.NullString `db:"delivery_address"`
		Region       sql.NullString `db:"delivery_region"`
		Email        sql.NullString `db:"delivery_email"`

		// Item
		ItemID          sql.NullInt64  `db:"item_id"`
		ChrtID          sql.NullInt64  `db:"item_chrt_id"`
		ItemTrackNumber sql.NullString `db:"item_track_number"`
		Price           sql.NullInt64  `db:"item_price"`
		Rid             sql.NullString `db:"item_rid"`
		Name            sql.NullString `db:"item_name"`
		Sale            sql.NullInt64  `db:"item_sale"`
		Size            sql.NullString `db:"item_size"`
		TotalPrice      sql.NullInt64  `db:"item_total_price"`
		NmID            sql.NullInt64  `db:"item_nm_id"`
		Brand           sql.NullString `db:"item_brand"`
		Status          sql.NullInt64  `db:"item_status"`
		ItemOrderID     sql.NullInt64  `db:"item_order_id"`
	}

	query := `SELECT
        o.id AS order_id,
        o.order_uid AS order_order_uid,
        o.track_number AS order_track_number,
        o.entry AS order_entry,
        o.locale AS order_locale,
        o.internal_signature AS order_internal_signature,
        o.customer_id AS order_customer_id,
        o.delivery_id AS order_delivery_id,
        o.payment_id AS order_payment_id,
        o.shardkey AS order_shardkey,
        o.sm_id AS order_sm_id,
        o.date_created AS order_date_created,
        o.oof_shard AS order_oof_shard,

        p.id AS payment_id,
        p.transaction AS payment_transaction,
        p.request_id AS payment_request_id,
        p.currency AS payment_currency,
        p.provider AS payment_provider,
        p.amount AS payment_amount,
        p.payment_dt AS payment_payment_dt,
        p.bank AS payment_bank,
        p.delivery_cost AS payment_delivery_cost,
        p.goods_total AS payment_goods_total,
        p.custom_fee AS payment_custom_fee,

        d.id AS delivery_id,
        d.name AS delivery_name,
        d.phone AS delivery_phone,
        d.zip AS delivery_zip,
        d.city AS delivery_city,
        d.address AS delivery_address,
        d.region AS delivery_region,
        d.email AS delivery_email,

        i.id AS item_id,
        i.chrt_id AS item_chrt_id,
        i.track_number AS item_track_number,
        i.price AS item_price,
        i.rid AS item_rid,
        i.name AS item_name,
        i.sale AS item_sale,
        i.size AS item_size,
        i.total_price AS item_total_price,
        i.nm_id AS item_nm_id,
        i.brand AS item_brand,
        i.status AS item_status,
        i.order_id AS item_order_id
    FROM "order" o
    LEFT JOIN payment p ON o.payment_id = p.id
    LEFT JOIN delivery d ON o.delivery_id = d.id
    LEFT JOIN item i ON i.order_id = o.id
    WHERE o.id = $1
    ORDER BY i.id;`

	var rows []orderRow
	err := repo.pool.SelectContext(ctx, &rows, query, orderID)

	lat := time.Since(start).Seconds()
	repo.metrics.DBQueryDuration.WithLabelValues("get_order_by_id", "order_service").Observe(lat)

	if err != nil {
		repo.metrics.DBQueryErrors.WithLabelValues("get_order_by_id", "order_service").Inc()
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("order not found")
	}

	r := rows[0]

	scanPayment := func(r orderRow) models.Payment {
		return models.Payment{
			ID:           int(r.PaymentID.Int64),
			Transaction:  r.Transaction.String,
			RequestID:    r.RequestID.String,
			Currency:     r.Currency.String,
			Provider:     r.Provider.String,
			Amount:       int(r.Amount.Int64),
			PaymentDt:    int(r.PaymentDt.Int64),
			Bank:         r.Bank.String,
			DeliveryCost: int(r.DeliveryCost.Int64),
			GoodsTotal:   int(r.GoodsTotal.Int64),
			CustomFee:    int(r.CustomFee.Int64),
			OrderID:      r.OrderPaymentID,
		}
	}

	scanDelivery := func(r orderRow) models.Delivery {
		return models.Delivery{
			ID:      int(r.DeliveryID.Int64),
			Name:    r.DeliveryName.String,
			Phone:   r.Phone.String,
			Zip:     r.Zip.String,
			City:    r.City.String,
			Address: r.Address.String,
			Region:  r.Region.String,
			Email:   r.Email.String,
			OrderID: r.OrderDeliveryID,
		}
	}

	scanItem := func(r orderRow) models.Item {
		return models.Item{
			ID:          int(r.ItemID.Int64),
			ChrtID:      int(r.ChrtID.Int64),
			TrackNumber: r.ItemTrackNumber.String,
			Price:       int(r.Price.Int64),
			Rid:         r.Rid.String,
			Name:        r.Name.String,
			Sale:        int(r.Sale.Int64),
			Size:        r.Size.String,
			TotalPrice:  int(r.TotalPrice.Int64),
			NmID:        int(r.NmID.Int64),
			Brand:       r.Brand.String,
			Status:      int(r.Status.Int64),
			OrderID:     int(r.ItemOrderID.Int64),
		}
	}

	order := &broker.OrderMessage{
		Order: models.Order{
			ID:                r.OrderID,
			OrderUID:          r.OrderUID,
			TrackNumber:       r.TrackNumber,
			Entry:             r.Entry,
			Locale:            r.Locale,
			InternalSignature: r.InternalSignature,
			CustomerID:        r.CustomerID,
			DeliveryID:        r.OrderDeliveryID,
			PaymentID:         r.OrderPaymentID,
			Shardkey:          r.Shardkey,
			SmID:              r.SmID,
			DateCreated:       r.DateCreated,
			OofShard:          r.OofShard,
		},
		Payment:  scanPayment(r),
		Delivery: scanDelivery(r),
		Items:    []models.Item{},
	}

	for _, r := range rows {
		if r.ItemID.Valid {
			order.Items = append(order.Items, scanItem(r))
		}
	}

	return order, nil
}

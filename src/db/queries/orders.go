package queries

import (
	"context"
	"fmt"
	"log"
	"orders/src/db/models"
)

func (service *DBService) CreateOrder(ctx context.Context, orderDto *models.Order) (models.Order, error) {
	var order models.Order

	query := `
    INSERT INTO "order" (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created,
	oof_shard, delivery_id, payment_id)
	VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created,
	:oof_shard, :delivery_id, :payment_id)
	RETURNING id, order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created,
	oof_shard, delivery_id, payment_id;
    `

	rows, err := service.DB.Pool.NamedQueryContext(ctx, query, &orderDto)

	if err != nil {
		return models.Order{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.StructScan(&order); err != nil {
			fmt.Printf("Error while parsing rows %v", err)
			return models.Order{}, err
		}
	}

	return order, nil
}

func (service *DBService) GetOrderById(ctx context.Context, orderId int) (models.Order, error) {
	var order models.Order

	query := `select *
			from "order"
			where id = $1;`

	err := service.DB.Pool.Get(&order, query, orderId)

	if err != nil {
		log.Fatalf("Error in GetOrderById: %v", err)
		return models.Order{}, err
	}

	return order, nil
}

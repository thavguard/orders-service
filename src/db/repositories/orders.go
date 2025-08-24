package repositories

import (
	"context"
	"fmt"
	"orders/src/db/models"
	"orders/src/internal/broker"
)

func (service *DBRepository) CreateOrder(ctx context.Context, orderDto *models.Order) (models.Order, error) {
	var order models.Order

	query := `
    INSERT INTO "order" (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created,
	oof_shard, delivery_id, payment_id)
	VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created,
	:oof_shard, :delivery_id, :payment_id)
	RETURNING id, order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created,
	oof_shard, delivery_id, payment_id;
    `

	rows, err := service.DB.Pool.NamedQueryContext(ctx, query, orderDto)

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

func (service *DBRepository) GetOrderById(ctx context.Context, orderId int) (*broker.OrderMessage, error) {
	var order *broker.OrderMessage

	var orderRaw models.Order
	query := `select * from "order" where id = $1;`

	err := service.DB.Pool.Get(&orderRaw, query, orderId)

	if err != nil {
		fmt.Printf("Error in GetOrderById Pool.Select: %v\n", err)
		return &broker.OrderMessage{}, err
	}

	delivery, err := service.GetDeliveryByOrderId(ctx, orderId)

	if err != nil {
		fmt.Printf("Error in GetDeliveryByOrderId: %v\n", err)
		return &broker.OrderMessage{}, err
	}

	payment, err := service.GetPaymentByOrderId(ctx, orderId)

	if err != nil {
		fmt.Printf("Error in GetPaymentByOrderId: %v\n", err)
		return &broker.OrderMessage{}, err
	}

	items, err := service.GetItemsByOrderId(ctx, orderId)

	if err != nil {
		fmt.Printf("Error in GetItemsByOrderId: %v\n", err)
		return &broker.OrderMessage{}, err
	}

	order = &broker.OrderMessage{
		Delivery: delivery,
		Payment:  payment,
		Items:    items,

		Order: &orderRaw,
	}

	fmt.Println(order)

	return order, nil
}

package queries

import (
	"context"
	"fmt"
	"log"
	"orders/src/db/models"
)

func (service *DBService) CreatePayment(ctx context.Context, paymentDto *models.Payment) (models.Payment, error) {
	var payment models.Payment

	query := `
     INSERT INTO payment (currency, delivery_cost, provider,
amount, payment_dt, bank, request_id, transaction, custom_fee, goods_total)
VALUES (:currency, :delivery_cost, :provider,
:amount, :payment_dt, :bank, :request_id, :transaction, :custom_fee, :goods_total)
RETURNING id, currency, delivery_cost, provider,
amount, payment_dt, bank, request_id, transaction, custom_fee, goods_total;
    `

	rows, err := service.DB.Pool.NamedQueryContext(ctx, query, &paymentDto)

	if err != nil {
		return models.Payment{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.StructScan(&payment); err != nil {
			fmt.Printf("Error while parsing rows %v", err)
			return models.Payment{}, err
		}
	}

	return payment, nil
}

func (service *DBService) GetPaymentById(ctx context.Context, paymentId int) (models.Payment, error) {
	var payment models.Payment

	query := `select *
			from payment
			where id = $1;`

	err := service.DB.Pool.Get(&payment, query, paymentId)

	if err != nil {
		log.Fatalf("Error in GetOrderById: %v", err)
		return models.Payment{}, err
	}

	return payment, nil
}

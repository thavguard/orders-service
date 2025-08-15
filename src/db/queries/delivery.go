package queries

import (
	"context"
	"fmt"
	"log"
	"orders/src/db/models"
)

func (service *DBService) CreateDelivery(ctx context.Context, deliveryDto *models.Delivery) (models.Delivery, error) {
	var delivery models.Delivery

	query := `
     INSERT INTO delivery (name, phone, zip, city, address, region, email)
VALUES (:name, :phone, :zip, :city, :address, :region, :email)
RETURNING id, name, phone, zip, city, address, region, email;
    `

	rows, err := service.DB.Pool.NamedQueryContext(ctx, query, &deliveryDto)

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

func (service *DBService) GetDeliveryById(ctx context.Context, deliveryId int) (models.Delivery, error) {
	var delivery models.Delivery

	query := `select *
			from delivery
			where id = $1;`

	err := service.DB.Pool.Get(&delivery, query, deliveryId)

	if err != nil {
		log.Fatalf("Error in GetDeliveryById: %v", err)
		return models.Delivery{}, err
	}

	return delivery, nil
}

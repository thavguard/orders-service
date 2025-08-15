package queries

import (
	"context"
	"fmt"
	"log"
	"orders/src/db/models"
)

func (service *DBService) CreateItem(ctx context.Context, itemDto *models.Item) (models.Item, error) {
	var item models.Item

	query := `
     INSERT INTO item (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id)
VALUES (:chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status, :order_id)
RETURNING id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id;
    `

	rows, err := service.DB.Pool.NamedQueryContext(ctx, query, &itemDto)

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

func (service *DBService) GetItemById(ctx context.Context, itemId int) (models.Item, error) {
	var item models.Item

	query := `select *
			from item
			where id = $1;`

	err := service.DB.Pool.Get(&item, query, itemId)

	if err != nil {
		log.Fatalf("Error in GetItemById: %v", err)
		return models.Item{}, err
	}

	return item, nil
}

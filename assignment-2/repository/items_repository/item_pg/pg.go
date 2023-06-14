package item_pg

import (
	"database/sql"
	"fmt"

	"assignment_2/models"
	"assignment_2/pkg/errs"
	"assignment_2/repository/items_repository"
)

type itemPG struct {
	db *sql.DB
}

func NewItemPG(db *sql.DB) items_repository.ItemRepository {
	return &itemPG{
		db: db,
	}
}


func (i *itemPG) generatePlaceholder(dataAmount int) string {
	start := "("

	for i := 1; i <= dataAmount; i++ {
		if i == dataAmount {
			start += fmt.Sprintf("$%d)", i+1)
			continue
		}
		start += fmt.Sprintf("$%d,", i+1)
	}

	return start
}

func (i *itemPG) findItemsByItemCodeQuery(dataAmount int, orderId int) string {
	query := `
	SELECT item_id, item_code, quantity, description, order_id, created_at, updated_at
	FROM "items"
	WHERE order_id=$1 AND item_code IN 
	`

	param := i.generatePlaceholder(dataAmount)

	return query + param
}

func (i *itemPG) FindItemsByItemCodes(itemCodes []string, orderId int) ([]*models.Items, errs.MessageErr) {
	query := i.findItemsByItemCodeQuery(len(itemCodes), orderId)

	args := []interface{}{orderId}
	for _, value := range itemCodes {
		args = append(args, value)
	}

	rows, err := i.db.Query(query, args...)
	if err != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	items := []*models.Items{}
	for rows.Next() {
		item := models.Items{}

		err := rows.Scan(&item.ItemId, &item.ItemCode, &item.Quantity, &item.Description, &item.OrderId, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, errs.NewInternalServerError("something went wrong")
		}

		items = append(items, &item)
	}

	return items, nil
}




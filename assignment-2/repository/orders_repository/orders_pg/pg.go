package orders_pg

import (
	"database/sql"
	"fmt"
	"time"
	"assignment_2/models"
	"assignment_2/pkg/errs"
	"assignment_2/repository/orders_repository"
)


const (

	createOrderQuery = `
		INSERT INTO "orders"
			(
				ordered_at,
				customer_name
			)
		VALUES($1, $2)
		RETURNING order_id, customer_name, ordered_at, created_at,updated_at
	`
	createItemQuery = `
		INSERT INTO "items"
			(
				item_code,
				description,
				quantity,
				order_id
			)
		VALUES($1, $2, $3, $4)
		RETURNING item_id, item_code, description, quantity, order_id, created_at, updated_at
	`

	getAllQueryOrder = `
		SELECT order_id, customer_name, ordered_at, created_at, updated_at 
		FROM "orders";
	`

	updateQueryOrder = `
		UPDATE "orders"
		SET ordered_at = $2,
		customer_name = $3,
		updated_at = $4
		WHERE order_id = $1
		RETURNING order_id, customer_name, ordered_at, created_at, updated_at
	`
	updateQueryItem = `
		UPDATE "items"
		SET description = $3,
		quantity = $4,
		updated_at = $5
		WHERE order_id = $1 AND item_code = $2
		RETURNING item_id, item_code, description, quantity, order_id, created_at, updated_at
	`
	getOrderQuery = `
		SELECT order_id, customer_name, ordered_at, created_at, updated_at
		FROM "orders"
		WHERE order_id=$1;
	`

	getItemsQuery = `
		SELECT item_id, item_code, description, quantity, order_id, created_at, updated_at
		FROM "items"
		WHERE order_id = $1;
		
	`

	deleteOrderQuery = `
		DELETE FROM "orders"
		WHERE order_id = $1;
	`

	selectOrderQuery = `
		SELECT * FROM "orders"
		WHERE order_id = $1
	`
)

type orderPG struct {
	db *sql.DB
}

func NewOrderPG(db *sql.DB) orders_repository.OrderRepository {
	return  &orderPG{db:db}  
}

func (o *orderPG) CreateOrder(orderPayload models.Order, itemsPayload []models.Items) (*models.Order, errs.MessageErr) {
	tx, err := o.db.Begin()
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to begin transaction")
	}

	orderRow := tx.QueryRow(createOrderQuery, orderPayload.OrderedAt, orderPayload.CustomerName)
	var order models.Order
	err = orderRow.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return nil, errs.NewInternalServerError("Failed to create new order")
	}

	items := []models.Items{}
	for _, item := range itemsPayload {
		itemRow := tx.QueryRow(createItemQuery, item.ItemCode, item.Description, item.Quantity, order.OrderId)
		var newItem models.Items
		err = itemRow.Scan(&newItem.ItemId, &newItem.ItemCode, &newItem.Description, &newItem.Quantity, &newItem.OrderId, &newItem.CreatedAt, &newItem.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return nil, errs.NewInternalServerError("Failed to create new item")
		}
		items = append(items, newItem)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, errs.NewInternalServerError("Failed to commit transaction")
	}

	order.Items = items

	return &order, nil
}

func (o *orderPG) GetAllOrder() ([]models.Order, errs.MessageErr) {
	rows, err := o.db.Query(getAllQueryOrder)
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to retrieve orders")
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, errs.NewInternalServerError("Failed to scan order row")
		}

		itemsRows, err := o.db.Query(getItemsQuery, order.OrderId)
		if err != nil {
			return nil, errs.NewInternalServerError("Failed to retrieve items")
		}
		defer itemsRows.Close()

		items := []models.Items{}
		for itemsRows.Next() {
			var item models.Items
			err := itemsRows.Scan(&item.ItemId, &item.ItemCode, &item.Description, &item.Quantity, &item.OrderId, &item.CreatedAt, &item.UpdatedAt)
			if err != nil {
				return nil, errs.NewInternalServerError("Failed to scan item row")
			}
			items = append(items, item)
		}

		order.Items = items
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *orderPG) GetOrderById(orderId int) (*models.Order, errs.MessageErr) {
	row := o.db.QueryRow(getOrderQuery, orderId)

	var order models.Order
	err := row.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFound("Order not found")
		} else {
			return nil, errs.NewInternalServerError("Something went wrong")
		}
	}

	itemsRows, err := o.db.Query(getItemsQuery, order.OrderId)
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to retrieve items")
	}
	defer itemsRows.Close()

	var items []models.Items
	for itemsRows.Next() {
		var item models.Items
		err := itemsRows.Scan(&item.ItemId, &item.ItemCode, &item.Description, &item.Quantity, &item.OrderId, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, errs.NewInternalServerError("Failed to scan item row")
		}
		items = append(items, item)
	}

	order.Items = items

	return &order, nil
}

func (o *orderPG) UpdateOrder(orderPayload models.Order, itemsPayload []models.Items) (*orders_repository.OrderItem, errs.MessageErr) {
	tx, err := o.db.Begin()
	if err != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	row := tx.QueryRow(updateQueryOrder, orderPayload.OrderId, orderPayload.OrderedAt, orderPayload.CustomerName, time.Now().Format(time.RFC3339))

	order := models.Order{}

	err = row.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return nil, errs.NewInternalServerError("something went wrong")
	}

	items := []models.Items{}
	if len(itemsPayload) != 0 {
		for _, eachItem := range itemsPayload {
			row = tx.QueryRow(updateQueryItem, order.OrderId, eachItem.ItemCode, eachItem.Description, eachItem.Quantity, time.Now().Format(time.RFC3339))

			item := models.Items{}

			err := row.Scan(&item.ItemId, &item.ItemCode, &item.Description, &item.Quantity, &item.OrderId, &item.CreatedAt, &item.UpdatedAt)
			if err != nil {
				tx.Rollback()
				return nil, errs.NewInternalServerError("something went wrong")
			}

			items = append(items, item)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, errs.NewInternalServerError("something went wrong")
	}

	result := &orders_repository.OrderItem{
		Order: order,
		Items: items,
	}

	return result, nil
}

func (o *orderPG) DeleteOrders(orderId int) errs.MessageErr {
	row := o.db.QueryRow(selectOrderQuery, orderId)

	var order models.Order
	if err := row.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.CreatedAt, &order.UpdatedAt); err != nil {
		fmt.Println(err)
		return errs.NewNotFound("Order Not Found")
	}

	_, err := o.db.Exec(deleteOrderQuery, orderId)
	if err != nil {
		return errs.NewInternalServerError("something went wrong")
	}
	return nil
}
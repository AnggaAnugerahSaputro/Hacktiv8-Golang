package orders_repository

import (
	"assignment_2/pkg/errs"
	"assignment_2/models"
)


type OrderRepository interface {
	CreateOrder(orderPayload models.Order, itemsPayload []models.Items) (*models.Order, errs.MessageErr)
	GetAllOrder() ([]models.Order, errs.MessageErr)
	GetOrderById(orderId int) (*models.Order, errs.MessageErr)
	UpdateOrder(models.Order, []models.Items) (*OrderItem, errs.MessageErr)
	DeleteOrders(orderId int) errs.MessageErr
}
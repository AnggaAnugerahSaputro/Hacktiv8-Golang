package services

import (
	"net/http"
	"fmt"
	"time"

	"assignment_2/dto"
	"assignment_2/models"
	"assignment_2/pkg/errs"
	"assignment_2/repository/orders_repository"
)

type orderService struct {
	orderRepo   orders_repository.OrderRepository
	itemService ItemService
}

type OrderService interface {
	CreateOrder(payload dto.NewOrderRequest) (*dto.NewOrderResponse, errs.MessageErr)
	UpdateOrder(orderId int, payload dto.NewOrderRequest) (*dto.NewOrderUpdateResponse, errs.MessageErr)
	GetOrderById(orderId int) (*dto.NewOrderUpdateResponse, errs.MessageErr)
	GetAllOrder() (*dto.GetAllOrderResponse, errs.MessageErr)
	DeleteOrders(orderId int) (*dto.DeleteOrderResponse, errs.MessageErr)
}

func NewOrderService(orderRepo orders_repository.OrderRepository, itemService ItemService) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		itemService: itemService,
	}
}

func (o *orderService) CreateOrder(payload dto.NewOrderRequest) (*dto.NewOrderResponse, errs.MessageErr) {
	orderPayload := o.buildCreateOrderPayload(payload)
	itemsPayload := o.buildCreateItemsPayload(payload.Items)

	newOrder, err := o.orderRepo.CreateOrder(orderPayload, itemsPayload)
	if err != nil {
		return nil, err
	}

	response := buildNewOrderResponse(newOrder)

	return response, nil
}

func (o *orderService) GetAllOrder() (*dto.GetAllOrderResponse, errs.MessageErr) {
	allOrder, err := o.orderRepo.GetAllOrder()
	if err != nil {
		return nil, err
	}
	response := buildAllOrderResponse(allOrder)

	return response, nil
}

func (o *orderService) GetOrderById(orderId int) (*dto.NewOrderUpdateResponse, errs.MessageErr) {
	order, err := o.orderRepo.GetOrderById(orderId)
	if err != nil {
		return nil, err
	}

	orderResponse := mapOrderToResponse(*order)

	response := &dto.NewOrderUpdateResponse{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("Get Order with ID %d Success", orderId),
		Data:       orderResponse,
	}

	return response, nil
}

func (o *orderService) UpdateOrder(orderId int, payload dto.NewOrderRequest) (*dto.NewOrderUpdateResponse, errs.MessageErr) {
	itemCodes := payload.ItemsToItemCode()

	itemsPayload := []models.Items{}
	if itemCodes != nil {
		_, err := o.itemService.FindItemsByItemCodes(itemCodes, orderId)
		if err != nil {
			return nil, err
		}
		itemsPayload = o.BuildItemsPayload(payload.Items)
	}

	orderPayload := o.BuildOrderPayload(orderId, payload.OrderedAt, payload.CustomerName)

	orderItem, err := o.orderRepo.UpdateOrder(orderPayload, itemsPayload)
	if err != nil {
		return nil, err
	}

	itemsResponse := []dto.ItemResponse{}
	if len(orderItem.Items) != 0 {
		for _, eachItem := range orderItem.Items {
			itemResponse := eachItem.ItemToItemResponse()

			itemsResponse = append(itemsResponse, itemResponse)
		}
	}

	result := buildNewOrderUpdateResponse(orderItem, itemsResponse)

	return result, nil
}

func (o *orderService) DeleteOrders(orderId int) (*dto.DeleteOrderResponse, errs.MessageErr) {
	_, err := o.orderRepo.GetOrderById(orderId)
	if err != nil {
		return nil, err
	}

	err = o.orderRepo.DeleteOrders(orderId)
	if err != nil {
		return nil, err
	}

	response := &dto.DeleteOrderResponse{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("Order with ID %d has been deleted", orderId),
	}

	return response, nil
}


// payload Create Order
func (o *orderService) buildCreateOrderPayload(payload dto.NewOrderRequest) models.Order {
	orderPayload := models.Order{
		OrderedAt:    payload.OrderedAt,
		CustomerName: payload.CustomerName,
	}

	return orderPayload
}

// payload Create Item Order
func (o *orderService) buildCreateItemsPayload(items []dto.NewItemRequest) []models.Items {
	itemsPayload := []models.Items{}
	for _, eachItem := range items {
		item := models.Items{
			ItemCode:    eachItem.ItemCode,
			Quantity:    eachItem.Quantity,
			Description: eachItem.Description,
		}
		itemsPayload = append(itemsPayload, item)
	}

	return itemsPayload
}

// payload Update  Order
func (o *orderService) BuildOrderPayload(orderId int, orderedAt time.Time, customerName string) models.Order {
	orderPayload := models.Order{
		OrderId:      orderId,
		OrderedAt:    orderedAt,
		CustomerName: customerName,
	}

	return orderPayload
}

// payload Update Item Order
func (o *orderService) BuildItemsPayload(items []dto.NewItemRequest) []models.Items {
	itemsPayload := []models.Items{}
	for _, eachItem := range items {
		item := models.Items{
			ItemCode:    eachItem.ItemCode,
			Description: eachItem.Description,
			Quantity:    eachItem.Quantity,
		}

		itemsPayload = append(itemsPayload, item)
	}

	return itemsPayload
}

// Response Create Order
func buildNewOrderResponse(order *models.Order) *dto.NewOrderResponse {
	var newItem []dto.NewItemRequest

	for _, eachItem := range order.Items {
		item := dto.NewItemRequest{
			ItemCode:    eachItem.ItemCode,
			Description: eachItem.Description,
			Quantity:    eachItem.Quantity,
		}

		newItem = append(newItem, item)
	}

	return &dto.NewOrderResponse{
		StatusCode: http.StatusCreated,
		Message: "Created Order successfully",
		Data: dto.NewOrderRequest{
			OrderedAt:    order.OrderedAt,
			CustomerName: order.CustomerName,
			Items:        newItem,
		},
	}
}

// Response Update
func buildNewOrderUpdateResponse(orderItem *orders_repository.OrderItem, itemsResponse []dto.ItemResponse) *dto.NewOrderUpdateResponse {
	result := &dto.NewOrderUpdateResponse{
		Message: "Update order successfully",
		Data: dto.OrderResponse{
			OrderId:      orderItem.Order.OrderId,
			OrderedAt:    orderItem.Order.OrderedAt,
			CustomerName: orderItem.Order.CustomerName,
			CreatedAt:    orderItem.Order.CreatedAt,
			UpdatedAt:    orderItem.Order.UpdatedAt,
			Items:        itemsResponse,
		},
		StatusCode: http.StatusOK,
	}

	return result
}

// Get Order
func mapOrderToResponse(order models.Order) dto.OrderResponse {
	var items []dto.ItemResponse
	for _, eachItem := range order.Items {
		item := dto.ItemResponse{
			Id:          eachItem.ItemId,
			ItemCode:    eachItem.ItemCode,
			Description: eachItem.Description,
			Quantity:    eachItem.Quantity,
			OrderId:     eachItem.OrderId,
		}

		items = append(items, item)
	}

	orderResponse := dto.OrderResponse{
		OrderId:      order.OrderId,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		OrderedAt:    order.OrderedAt,
		CustomerName: order.CustomerName,
		Items:        items,
	}

	return orderResponse
}

// Get All Order
func buildAllOrderResponse(allOrder []models.Order) *dto.GetAllOrderResponse {
	var orders []dto.OrderResponse
	for _, order := range allOrder {
		var items []dto.ItemResponse
		for _, item := range order.Items {
			i := dto.ItemResponse{
				Id:          item.ItemId,
				ItemCode:    item.ItemCode,
				Description: item.Description,
				Quantity:    item.Quantity,
				OrderId:     item.OrderId,
				CreatedAt:   item.CreatedAt,
				UpdatedAt:   item.UpdatedAt,
			}

			items = append(items, i)
		}

		o := dto.OrderResponse{
			OrderId:      order.OrderId,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.CreatedAt,
			CustomerName: order.CustomerName,
			Items:        items,
		}

		orders = append(orders, o)
	}

	response := &dto.GetAllOrderResponse{
		StatusCode: http.StatusOK,
		Message:    "Get All Orders Successfully",
		Data:       orders,
	}

	return response
}

package dto

import (
	"time"
)

type NewOrderRequest struct {
	OrderedAt    time.Time        `json:"orderedAt" binding:"required"`
	CustomerName string           `json:"customerName" binding:"required"`
	Items        []NewItemRequest `json:"items" binding:"required"`
}


func (o *NewOrderRequest) ItemsToItemCode() []string {
	itemCodes := []string{}

	for _, value := range o.Items {
		itemCodes = append(itemCodes, value.ItemCode)
	}

	return itemCodes
}

type NewOrderResponse struct {
	StatusCode int           	`json:"code"`
	Message    string        	`json:"message"`
	Data       NewOrderRequest 	`json:"data"`
}

type NewOrderUpdateResponse struct {
	StatusCode int           	`json:"code"`
	Message    string        	`json:"message"`
	Data       OrderResponse 	`json:"data"`
}

type OrderResponse struct {
	OrderId      int            `json:"id"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	OrderedAt    time.Time      `json:"orderedAt"`
	CustomerName string         `json:"customerName"`
	Items        []ItemResponse `json:"items"`
}

type GetAllOrderResponse struct {
	StatusCode int             `json:"code"`
	Message    string          `json:"message"`
	Data       []OrderResponse `json:"data"`
}

type DeleteOrderResponse struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

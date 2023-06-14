package dto

import (
	"time"
)

type NewItemRequest struct {
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}
type ItemResponse struct {
	Id          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ItemCode    string    `json:"itemcode"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	OrderId     int       `json:"orderid"`
}
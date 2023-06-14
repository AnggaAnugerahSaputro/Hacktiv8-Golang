package models

import (
	"time"
	"assignment_2/dto"
)

type Items struct {
	ItemId      int
	ItemCode    string
	Description string
	Quantity    int
	OrderId     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (i *Items) ItemToItemResponse() dto.ItemResponse {
	return dto.ItemResponse{
		Id:          i.ItemId,
		ItemCode:    i.ItemCode,
		Quantity:    i.Quantity,
		Description: i.Description,
		OrderId:     i.OrderId,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}
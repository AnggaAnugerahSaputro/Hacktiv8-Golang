package models

import "time"

type Order struct {
	OrderId      int
	CustomerName string
	OrderedAt    time.Time
	Items		 []Items
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

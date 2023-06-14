package orders_repository


import "assignment_2/models"


type OrderItem struct {
	Message    string          `json:"message"`
	StatusCode int             `json:"code"`
	Order models.Order
	Items []models.Items
}
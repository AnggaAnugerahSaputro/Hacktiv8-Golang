package items_repository

import (
	"assignment_2/pkg/errs"
	"assignment_2/models"
)

type ItemRepository interface {
	FindItemsByItemCodes([]string, int) ([]*models.Items, errs.MessageErr)
}
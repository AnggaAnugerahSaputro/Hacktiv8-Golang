package services

import (
	"fmt"

	"assignment_2/models"
	"assignment_2/pkg/errs"
	"assignment_2/repository/items_repository"
)

type itemService struct {
	itemRepo items_repository.ItemRepository
}

type ItemService interface {
	FindItemsByItemCodes([]string, int) ([]*models.Items, errs.MessageErr)
}

func NewItemService(itemRepo items_repository.ItemRepository) ItemService {
	return &itemService{
		itemRepo: itemRepo,
	}
}

func (i *itemService) FindItemsByItemCodes(itemCodes []string, orderId int) ([]*models.Items, errs.MessageErr) {
	item, err := i.itemRepo.FindItemsByItemCodes(itemCodes, orderId)

	if err != nil {
		return nil, err
	}

	for _, eachItemCode := range itemCodes {
		isFound := false

		for _, eachItem := range item {
			if eachItemCode == eachItem.ItemCode {
				isFound = true
				break
			}
		}

		if !isFound {
			return nil, errs.NewNotFound(fmt.Sprintf("item with code %s does not exist", eachItemCode))
		}
	}

	return item, err
}

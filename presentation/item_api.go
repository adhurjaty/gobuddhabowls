package presentation

import (
	"buddhabowls/models"
	"errors"
)

// ItemAPI is an object for serving order and vendor items to UI
type ItemAPI struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Category      CategoryAPI `json:"Category"`
	Index         int         `json:"index"`
	Count         float64     `json:"count,string,omitempty"`
	Price         float64     `json:"price,string,omitempty"`
	PurchasedUnit string      `json:"purchased_unit,string,omitempty"`
}

type ItemsAPI []ItemAPI

// NewItemAPI converts an order/vendor/inventory item to an api item
// TODO: create interface for these ^ types of items
func NewItemAPI(m interface{}) ItemAPI {
	item := ItemAPI{}
	switch m.(type) {
	case models.OrderItem:
		orderItem, _ := m.(models.OrderItem)
		return ItemAPI{
			ID:       orderItem.ID.String(),
			Name:     orderItem.InventoryItem.Name,
			Category: NewCategoryAPI(orderItem.InventoryItem.Category),
			Count:    orderItem.Count,
			Price:    orderItem.Price,
		}

	case models.VendorItem:
		vendorItem, _ := m.(models.VendorItem)
		return ItemAPI{
			ID:            vendorItem.ID.String(),
			Name:          vendorItem.InventoryItem.Name,
			Category:      NewCategoryAPI(vendorItem.InventoryItem.Category),
			Price:         vendorItem.Price,
			PurchasedUnit: vendorItem.PurchasedUnit.String,
		}
	case models.InventoryItem:
		invItem, _ := m.(models.InventoryItem)
		return ItemAPI{
			ID:       invItem.ID.String(),
			Name:     invItem.Name,
			Category: NewCategoryAPI(invItem.Category),
		}
	default:
		errors.New("Must supply OrderItem, VendorItem or InventoryItem type")
	}

	return item
}

// NewItemsAPI converts an order/vendor/inventory item slice to an api item slice
func NewItemsAPI(modelItems interface{}) ItemsAPI {
	modelSlice := modelItems.([]interface{})
	apis := make([]ItemAPI, len(modelSlice))
	for i, modelItem := range modelSlice {
		apis[i] = NewItemAPI(modelItem)
	}

	return apis
}

package presentation

import (
	"buddhabowls/models"
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
func NewItemAPI(item models.GenericItem) ItemAPI {
	itemAPI := ItemAPI{
		ID:       item.GetID().String(),
		Name:     item.GetName(),
		Category: NewCategoryAPI(item.GetCategory()),
		Index:    item.GetIndex(),
	}

	switch item.(type) {
	case models.OrderItem:
		orderItem, _ := item.(models.OrderItem)
		itemAPI.Count = orderItem.Count
		itemAPI.Price = orderItem.Price

	case models.VendorItem:
		vendorItem, _ := item.(models.VendorItem)
		itemAPI.Price = vendorItem.Price
		itemAPI.PurchasedUnit = vendorItem.PurchasedUnit.String
	}

	return itemAPI
}

// NewItemsAPI converts an order/vendor/inventory item slice to an api item slice
func NewItemsAPI(modelItems interface{}) ItemsAPI {
	var apis []ItemAPI
	// TODO: gotta be a better way to do this
	switch modelItems.(type) {
	case models.OrderItems:
		modelSlice := modelItems.(models.OrderItems)
		apis = make([]ItemAPI, len(modelSlice))
		for i, modelItem := range modelSlice {
			apis[i] = NewItemAPI(modelItem)
		}
	case models.VendorItems:
		modelSlice := modelItems.(models.VendorItems)
		apis = make([]ItemAPI, len(modelSlice))
		for i, modelItem := range modelSlice {
			apis[i] = NewItemAPI(modelItem)
		}
	case models.InventoryItems:
		modelSlice := modelItems.(models.InventoryItems)
		apis = make([]ItemAPI, len(modelSlice))
		for i, modelItem := range modelSlice {
			apis[i] = NewItemAPI(modelItem)
		}
	}

	return apis
}

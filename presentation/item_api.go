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

// ConvertToAPI converts an order/vendor/inventory item to an api item
func (item *ItemAPI) ConvertToAPI(m interface{}) error {
	switch m.(type) {
	case models.OrderItem:
		orderItem, _ := m.(models.OrderItem)
		item.ID = orderItem.ID.String()
		item.Name = orderItem.InventoryItem.Name
		item.Category = CategoryAPI{}
		item.Category.ConvertToAPI(orderItem.GetCategory())
		item.Count = orderItem.Count
		item.Price = orderItem.Price
	case models.VendorItem:
		vendorItem, _ := m.(models.VendorItem)
		item.ID = vendorItem.ID.String()
		item.Name = vendorItem.InventoryItem.Name
		item.Category = CategoryAPI{}
		item.Category.ConvertToAPI(vendorItem.GetCategory())
		item.Price = vendorItem.Price
		item.PurchasedUnit = vendorItem.PurchasedUnit.String
	case models.InventoryItem:
		invItem, _ := m.(models.InventoryItem)
		item.ID = invItem.ID.String()
		item.Name = invItem.Name
		item.Category = CategoryAPI{}
		item.Category.ConvertToAPI(invItem.Category)
	default:
		errors.New("Must supply OrderItem, VendorItem or InventoryItem type")
	}

	return nil
}

// ConvertToAPI converts an order/vendor/inventory item slice to an api item slice
func (items *ItemsAPI) ConvertToAPI(m interface{}) error {
	modelItems, ok := m.(models.InventoryItems)
	if !ok {
		return errors.New("Must supply OrderItem, VendorItem or InventoryItem type")
	}

	apis := ItemsAPI{}
	for _, modelItem := range modelItems {
		api := ItemAPI{}
		if err := api.ConvertToAPI(modelItem); err != nil {
			return err
		}
		apis = append(apis, api)
	}

	items = &apis
	return nil
}

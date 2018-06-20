package presentation

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop/nulls"
)

// PurchaseOrderAPI purchase order information to pass to the UI
type PurchaseOrderAPI struct {
	ID           string     `json:"id"`
	Vendor       VendorAPI  `json:"Vendor"`
	OrderDate    nulls.Time `json:"order_date"`
	ReceivedDate nulls.Time `json:"received_date,time,omitempty"`
	ShippingCost float64    `json:"shipping_cost" db:"shipping_cost"`
	Items        ItemsAPI   `json:"Items"`
}

type PurchaseOrdersAPI []PurchaseOrderAPI

// NewPurchaseOrderAPI converts a purchase order to an api purchase order
func NewPurchaseOrderAPI(purchaseOrder *models.PurchaseOrder) PurchaseOrderAPI {
	return PurchaseOrderAPI{
		ID:           purchaseOrder.ID.String(),
		Vendor:       NewVendorAPI(&purchaseOrder.Vendor),
		OrderDate:    purchaseOrder.OrderDate,
		ReceivedDate: purchaseOrder.ReceivedDate,
		ShippingCost: purchaseOrder.ShippingCost,
		Items:        NewItemsAPI(purchaseOrder.Items),
	}
}

// NewPurchaseOrdersAPI converts a purchase order slice to an api purchase order slice
func NewPurchaseOrdersAPI(purchaseOrders *models.PurchaseOrders) PurchaseOrdersAPI {
	apis := make([]PurchaseOrderAPI, len(*purchaseOrders))
	for i, po := range *purchaseOrders {
		apis[i] = NewPurchaseOrderAPI(&po)
	}

	return apis
}

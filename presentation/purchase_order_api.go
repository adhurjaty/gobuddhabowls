package presentation

import (
	"buddhabowls/models"
	"errors"
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

// ConvertToAPI converts a purchase order to an api purchase order
func (p *PurchaseOrderAPI) ConvertToAPI(m interface{}) error {
	purchaseOrder, ok := m.(models.PurchaseOrder)
	if !ok {
		return errors.New("Must supply PurchaseOrder type")
	}

	p.ID = purchaseOrder.ID.String()
	p.Vendor = VendorAPI{}
	p.Vendor.ConvertToAPI(purchaseOrder.Vendor)
	p.OrderDate = purchaseOrder.OrderDate
	p.ReceivedDate = purchaseOrder.ReceivedDate
	p.ShippingCost = purchaseOrder.ShippingCost
	p.Items = ItemsAPI{}
	p.Items.ConvertToAPI(purchaseOrder.Items)

	return nil
}

// ConvertToAPI converts a purchase order slice to an api purchase order slice
func (p *PurchaseOrdersAPI) ConvertToAPI(m interface{}) error {
	purchaseOrders, ok := m.(models.PurchaseOrders)
	if !ok {
		return errors.New("Must supply PurchaseOrders type")
	}

	apis := PurchaseOrdersAPI{}
	for _, po := range purchaseOrders {
		api := PurchaseOrderAPI{}
		if err := api.ConvertToAPI(po); err != nil {
			return err
		}
		apis = append(apis, api)
	}

	p = &apis
	return nil
}

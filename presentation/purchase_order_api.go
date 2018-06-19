package presentation

import (
	"github.com/gobuffalo/pop/nulls"
)

// PurchaseOrderAPI purchase order information to pass to the UI
type PurchaseOrderAPI struct {
	ID           string     `json:"id"`
	Vendor       VendorAPI  `json:"Vendor"`
	OrderDate    nulls.Time `json:"order_date"`
	ReceivedDate nulls.Time `json:"received_date"`
	ShippingCost float64    `json:"shipping_cost" db:"shipping_cost"`
	Items        []ItemAPI  `json:"Items"`
}

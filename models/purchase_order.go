package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type PurchaseOrder struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	VendorID     uuid.UUID  `json:"vendor_id" db:"vendor_id"`
	Vendor       Vendor     `belongs_to:"vendors" db:"-"`
	OrderDate    time.Time  `json:"order_date" db:"order_date"`
	ReceivedDate time.Time  `json:"received_date" db:"received_date"`
	ShippingCost float64    `json:"shipping_cost" db:"shipping_cost"`
	Items        OrderItems `has_many:"order_items" db:"-" fk_id:"order_id"`
}

// String is not required by pop and may be deleted
func (p PurchaseOrder) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PurchaseOrders is not required by pop and may be deleted
type PurchaseOrders []PurchaseOrder

// String is not required by pop and may be deleted
func (p PurchaseOrders) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PurchaseOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: p.OrderDate, Name: "OrderDate"},
		&validators.TimeIsPresent{Field: p.ReceivedDate, Name: "ReceivedDate"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PurchaseOrder) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PurchaseOrder) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// GetCost gets the total cost of the purchase order
func (p PurchaseOrder) GetCost() float64 {
	var cost float64
	for _, item := range p.Items {
		cost += item.Price
	}

	return cost + p.ShippingCost
}

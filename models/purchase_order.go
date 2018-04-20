package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/validate/validators"
	"time"

	"buddhabowls/helpers"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type PurchaseOrder struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	VendorID     uuid.UUID  `json:"vendor_id" db:"vendor_id"`
	Vendor       Vendor     `belongs_to:"vendors" db:"-"`
	OrderDate    nulls.Time `json:"order_date" db:"order_date"`
	ReceivedDate nulls.Time `json:"received_date" db:"received_date"`
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
		validate.ValidatorFunc(func(errors *validate.Errors) {
			if !p.OrderDate.Valid {
				errors.Add(validators.GenerateKey(p.Vendor.Name), "Must have order date")
			}
		}),
		validate.ValidatorFunc(func(errors *validate.Errors) {
			if p.OrderDate.Valid && p.ReceivedDate.Valid && p.OrderDate.Time.Unix() > p.ReceivedDate.Time.Unix() {
				errors.Add(validators.GenerateKey(p.Vendor.Name+" "+p.OrderDate.Time.String()),
					"Received date must be after order date")
			}
		}),
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
		cost += item.Price * item.Count
	}

	return cost + p.ShippingCost
}

// GetCategoryCosts gets a map of category -> cost map for the order
func (p PurchaseOrder) GetCategoryCosts() map[string]float64 {
	catCosts := map[string]float64{}

	for _, item := range p.Items {
		// fmt.Println(item)
		catCosts = helpers.AddToMap(catCosts, item.InventoryItem.Category, item.Price*item.Count)
	}

	return catCosts
}

// func (p *PurchaseOrder) GetYears() []int {

// }

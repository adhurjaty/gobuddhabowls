package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type OrderItem struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	InventoryItemID uuid.UUID     `json:"inventory_item_id" db:"inventory_item_id"`
	InventoryItem   InventoryItem `belongs_to:"inventory_items" db:"-"`
	OrderID         uuid.UUID     `json:"order_id" db:"order_id"`
	Order           PurchaseOrder `belongs_to:"purchase_orders" db:"-"`
	Price           float64       `json:"price" db:"price"`
	Count           float64       `json:"count" db:"count"`
}

// String is not required by pop and may be deleted
func (o OrderItem) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// OrderItems is not required by pop and may be deleted
type OrderItems []OrderItem

// String is not required by pop and may be deleted
func (o OrderItems) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *OrderItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *OrderItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *OrderItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

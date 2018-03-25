package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type VendorItem struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	InventoryItemID uuid.UUID     `json:"inventory_item_id" db:"inventory_item_id"`
	InventoryItem   InventoryItem `belongs_to:"inventory_item" db:"-"`
	VendorID        uuid.UUID     `json:"vendor_id" db:"vendor_id"`
	Vendor          Vendor        `belongs_to:"vendors" db:"-"`
	PurchasedUnit   string        `json:"purchased_unit" db:"purchased_unit"`
	Conversion      float64       `json:"conversion" db:"conversion"`
	Price           float64       `json:"price" db:"price"`
}

// String is not required by pop and may be deleted
func (v VendorItem) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// VendorItems is not required by pop and may be deleted
type VendorItems []VendorItem

// String is not required by pop and may be deleted
func (v VendorItems) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (v *VendorItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: v.PurchasedUnit, Name: "PurchasedUnit"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (v *VendorItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (v *VendorItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

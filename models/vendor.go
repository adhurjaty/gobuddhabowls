package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Vendor struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	Name         string      `json:"name" db:"name"`
	Email        string      `json:"email" db:"email"`
	PhoneNumber  string      `json:"phone_number" db:"phone_number"`
	Contact      string      `json:"contact" db:"contact"`
	ShippingCost float64     `json:"shipping_cost" db:"shipping_cost"`
	Items        VendorItems `has_many:"vendor_items" db:"-"`
}

// String is not required by pop and may be deleted
func (v Vendor) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Vendors is not required by pop and may be deleted
type Vendors []Vendor

// String is not required by pop and may be deleted
func (v Vendors) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (v *Vendor) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: v.Name, Name: "Name"},
		&validators.StringIsPresent{Field: v.Email, Name: "Email"},
		&validators.StringIsPresent{Field: v.PhoneNumber, Name: "PhoneNumber"},
		&validators.StringIsPresent{Field: v.Contact, Name: "Contact"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (v *Vendor) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (v *Vendor) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

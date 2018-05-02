package models

import (
	"encoding/json"
	"fmt"
	"time"

	"database/sql"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Vendor struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
	Name         string         `json:"name" db:"name"`
	Email        sql.NullString `json:"email" db:"email"`
	PhoneNumber  sql.NullString `json:"phone_number" db:"phone_number"`
	Contact      sql.NullString `json:"contact" db:"contact"`
	ShippingCost float64        `json:"shipping_cost" db:"shipping_cost"`
	Items        VendorItems    `has_many:"vendor_items" db:"-"`
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

// GetCategoryGroups gets all items the vendor sells grouped by category
func (v Vendor) GetCategoryGroups() map[InventoryItemCategory]VendorItems {
	outMap := make(map[InventoryItemCategory]VendorItems)

	fmt.Println(v)

	for _, item := range v.Items {
		itemList, ok := outMap[item.InventoryItem.Category]
		if ok {
			outMap[item.InventoryItem.Category] = append(itemList, item)
		} else {
			outMap[item.InventoryItem.Category] = VendorItems{item}
		}
		fmt.Println(outMap)
	}
	return outMap
}

// SelectValue returns the ID for select input tags
func (v Vendor) SelectValue() interface{} {
	return v.ID
}

// SelectLabel returs the name for select input tags
func (v Vendor) SelectLabel() string {
	return v.Name
}

// LoadVendor gets a vendor given ID
func LoadVendor(tx *pop.Connection, id string) (Vendor, error) {
	vendor := Vendor{}
	if err := tx.Eager().Find(&vendor, id); err != nil {
		return vendor, err
	}

	for i := 0; i < len(vendor.Items); i++ {
		if err := tx.Eager().Find(&vendor.Items[i], vendor.Items[i].ID); err != nil {
			return vendor, err
		}
		if err := tx.Eager().Find(&vendor.Items[i].InventoryItem, vendor.Items[i].InventoryItemID); err != nil {
			return vendor, err
		}
	}

	vendor.Items.Sort()

	return vendor, nil
}

// LoadVendors loads vendors and sub-models
func LoadVendors(q *pop.Query) (Vendors, error) {
	vendList := Vendors{}

	if err := q.All(&vendList); err != nil {
		return nil, err
	}

	// I don't love the fact that I need to load the nested models manually
	// TODO: look for a solution to eager loading nested objects
	for _, v := range vendList {
		for i := 0; i < len(v.Items); i++ {
			if err := q.Connection.Eager().Find(&v.Items[i], v.Items[i].ID); err != nil {
				return nil, err
			}
			if err := q.Connection.Eager().Find(&v.Items[i].InventoryItem, v.Items[i].InventoryItemID); err != nil {
				return nil, err
			}
		}

		v.Items.Sort()
	}

	return vendList, nil
}

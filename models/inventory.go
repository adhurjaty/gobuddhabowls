package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	// "github.com/gobuffalo/validate/validators"
)

type Inventory struct {
	ID        uuid.UUID           `json:"id" db:"id"`
	Date      time.Time           `json:"date" db:"date"`
	CreatedAt time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" db:"updated_at"`
	Items     CountInventoryItems `has_many:"count_inventory_items" db:"-"`
}

// String is not required by pop and may be deleted
func (i Inventory) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Inventories is not required by pop and may be deleted
type Inventories []Inventory

// String is not required by pop and may be deleted
func (i Inventories) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
// func (i *Inventory) Validate(tx *pop.Connection) (*validate.Errors, error) {
// 	return validate.Validate(
// 		&validators.TimeIsPresent{Field: i.Date, Name: "Date"},
// 	), nil
// }

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *Inventory) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *Inventory) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (i *Inventory) GetID() uuid.UUID {
	return i.ID
}

func (i *Inventory) GetItems() CompoundItems {
	return &i.Items
}

func (i *Inventories) ToModels() *[]Model {
	models := make([]Model, len(*i))
	for idx := range *i {
		models[idx] = &(*i)[idx]
	}

	return &models
}

func (i *Inventories) ToCompoundModels() *[]CompoundModel {
	models := make([]CompoundModel, len(*i))
	for idx := range *i {
		models[idx] = &(*i)[idx]
	}

	return &models
}

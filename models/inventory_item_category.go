package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type InventoryItemCategory struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Name       string    `json:"name" db:"name"`
	Background string    `json:"background" db:"background"`
	Index      int       `json:"index" db:"index"`
}

// String is not required by pop and may be deleted
func (i InventoryItemCategory) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// InventoryItemCategories is not required by pop and may be deleted
type InventoryItemCategories []InventoryItemCategory

// String is not required by pop and may be deleted
func (i InventoryItemCategories) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *InventoryItemCategory) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Name, Name: "Name"},
		&validators.StringIsPresent{Field: i.Background, Name: "Background"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *InventoryItemCategory) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *InventoryItemCategory) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (c InventoryItemCategory) GetID() uuid.UUID {
	return c.ID
}

func (c InventoryItemCategory) GetName() string {
	return c.Name
}

func (c InventoryItemCategory) GetBackground() string {
	return c.Background
}

func (c InventoryItemCategory) GetIndex() int {
	return c.Index
}

package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type PrepItem struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	InventoryItemID uuid.UUID     `json:"inventory_item_id" db:"inventory_item_id"`
	InvenotryItem   InventoryItem `belongs_to:"inventory_items" db:"-"`
	BatchRecipeID   uuid.UUID     `json:"batch_recipe_id" db:"batch_recipe_id"`
	BatchRecipe     Recipe        `belongs_to:"recipes" db:"-"`
	Conversion      float64       `json:"conversion" db:"conversion"`
	CountUnit       string        `json:"count_unit" db:"count_unit"`
	Index           int           `json:"index" db:"index"`
}

// String is not required by pop and may be deleted
func (p PrepItem) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PrepItems is not required by pop and may be deleted
type PrepItems []PrepItem

// String is not required by pop and may be deleted
func (p PrepItems) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PrepItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.CountUnit, Name: "CountUnit"},
		&validators.IntIsPresent{Field: p.Index, Name: "Index"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PrepItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PrepItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

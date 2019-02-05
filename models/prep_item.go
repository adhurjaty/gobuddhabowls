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
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	InventoryID uuid.UUID `json:"inventory_id" db:"inventory_id"`
	Inventory   Inventory `belongs_to:"inventories" db:"-"`
	Count       float64   `json:"count" db:"count"`
	RecipeID    uuid.UUID `json:"recipe_id" db:"recipe_id"`
	Recipe      Recipe    `belongs_to:"recipes" db:"-"`
	// Conversion is the number of recipe units in a prep item count
	Conversion float64 `json:"conversion" db:"conversion"`
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
		&validators.FuncValidator{
			Field: "Conversion",
			Name:  "Conversion",
			Fn: func() bool {
				return p.Conversion > 0
			},
		},
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

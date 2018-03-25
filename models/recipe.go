package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Recipe struct {
	ID                   uuid.UUID   `json:"id" db:"id"`
	CreatedAt            time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at" db:"updated_at"`
	Name                 string      `json:"name" db:"name"`
	RecipeUnit           string      `json:"recipe_unit" db:"recipe_unit"`
	RecipeUnitConversion string      `json:"recipe_unit_conversion" db:"recipe_unit_conversion"`
	Category             string      `json:"category" db:"category"`
	IsBatch              bool        `json:"is_batch" db:"is_batch"`
	Index                int         `json:"index" db:"index"`
	Items                RecipeItems `has_many:"recipe_items" db:"-"`
}

// String is not required by pop and may be deleted
func (r Recipe) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Recipes is not required by pop and may be deleted
type Recipes []Recipe

// String is not required by pop and may be deleted
func (r Recipes) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Recipe) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
		&validators.StringIsPresent{Field: r.RecipeUnit, Name: "RecipeUnit"},
		&validators.StringIsPresent{Field: r.RecipeUnitConversion, Name: "RecipeUnitConversion"},
		&validators.StringIsPresent{Field: r.Category, Name: "Category"},
		&validators.IntIsPresent{Field: r.Index, Name: "Index"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Recipe) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Recipe) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

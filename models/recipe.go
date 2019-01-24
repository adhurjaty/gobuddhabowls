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
	ID         uuid.UUID      `json:"id" db:"id"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	Name       string         `json:"name" db:"name"`
	Category   RecipeCategory `belongs_to:"recipe_categories" db:"-"`
	CategoryID uuid.UUID      `json:"category_id" db:"recipe_category_id"`
	RecipeUnit string         `json:"recipe_unit" db:"recipe_unit"`
	// RecipeUnitConversion is the number of recipe units in a yield
	//of a recipe
	RecipeUnitConversion float64     `json:"recipe_unit_conversion" db:"recipe_unit_conversion"`
	Items                RecipeItems `has_many:"recipe_items" db:"-"`
	Index                int         `json:"index" db:"index"`
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

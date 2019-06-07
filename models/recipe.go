package models

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Recipe struct {
	ID         uuid.UUID    `json:"id" db:"id"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" db:"updated_at"`
	Name       string       `json:"name" db:"name"`
	Category   ItemCategory `belongs_to:"item_categories" db:"-"`
	CategoryID uuid.UUID    `json:"category_id" db:"category_id"`
	RecipeUnit string       `json:"recipe_unit" db:"recipe_unit"`
	IsBatch    bool         `json:"is_batch" db:"is_batch"`
	// RecipeUnitConversion is the number of recipe units in a yield
	// of a recipe
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
		&validators.IntIsGreaterThan{
			Field:    r.Index,
			Name:     "Index",
			Compared: -1,
		},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Recipe) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return r.validateUniqueName(tx.Q())
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Recipe) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	query := tx.Where("id != ?", r.ID.String()).
		Where("is_batch = ?", r.IsBatch)
	return r.validateUniqueName(query)
}

func (r *Recipe) validateUniqueName(query *pop.Query) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	items := &Recipes{}
	if err := query.All(items); err != nil {
		return verrs, err
	}

	for _, item := range *items {
		if strings.ToLower(r.Name) == strings.ToLower(item.Name) {
			verrs.Add("Name", "Recipe name already exists")
			break
		}
	}

	return verrs, nil
}

func (r *Recipe) GetID() uuid.UUID {
	return r.ID
}
func (r *Recipe) GetName() string {
	return r.Name
}
func (r *Recipe) GetCategory() ItemCategory {
	return r.Category
}
func (r *Recipe) SetCategory(category ItemCategory) {
	r.Category = category
}
func (r *Recipe) GetCountUnit() string {
	return fmt.Sprintf("%d x %s", r.RecipeUnitConversion, r.RecipeUnit)
}
func (r *Recipe) GetIndex() int {
	return r.Index
}

func (r *Recipe) GetItems() CompoundItems {
	return &r.Items
}

func (r *Recipe) SetItems(items *[]CompoundItem) {
	r.Items = RecipeItems{}
	for _, item := range *items {
		r.Items = append(r.Items, *item.(*RecipeItem))
	}
}

func (r Recipe) GetSortValue() int {
	return r.Category.Index*1000 + r.Index
}

func (r *Recipes) Sort() {
	sort.Slice(*r, func(i, j int) bool {
		return (*r)[i].GetSortValue() < (*r)[j].GetSortValue()
	})
}

func (r *Recipes) ToGenericItems() *[]GenericItem {
	items := make([]GenericItem, len(*r))
	for i := 0; i < len(*r); i++ {
		items[i] = &(*r)[i]
	}

	return &items
}

func (r *Recipes) ToModels() *[]Model {
	models := make([]Model, len(*r))
	for idx := range *r {
		models[idx] = &(*r)[idx]
	}

	return &models
}

func (r *Recipes) ToCompoundModels() *[]CompoundModel {
	models := make([]CompoundModel, len(*r))
	for idx := range *r {
		models[idx] = &(*r)[idx]
	}

	return &models
}

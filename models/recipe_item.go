package models

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type RecipeItem struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	RecipeID        uuid.UUID     `json:"recipe_id" db:"recipe_id"`
	Recipe          Recipe        `belongs_to:"recipes" db:"-"`
	InventoryItemID uuid.NullUUID `json:"inventory_item_id" db:"inventory_item_id"`
	InventoryItem   InventoryItem `belongs_to:"inventory_items" db:"-"`
	BatchRecipeID   uuid.NullUUID `json:"batch_recipe_id" db:"batch_recipe_id"`
	BatchRecipe     Recipe        `belongs_to:"recipes" db:"-"`
	Measure         string        `json:"measure" db:"measure"`
	Count           float64       `json:"count" db:"count"`
}

// String is not required by pop and may be deleted
func (r RecipeItem) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// RecipeItems is not required by pop and may be deleted
type RecipeItems []RecipeItem

// String is not required by pop and may be deleted
func (r RecipeItems) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *RecipeItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{
			Field: r.Measure,
			Name:  "Measure",
		},
		&validators.FuncValidator{
			Field:   "BatchRecipeID",
			Name:    "BatchRecipeID",
			Message: "Recipe cannot contain itself... moron",
			Fn: func() bool {
				if !r.BatchRecipeID.Valid {
					return true
				}
				return r.ID.String() != r.BatchRecipeID.UUID.String()
			},
		},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *RecipeItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *RecipeItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (r *RecipeItem) GetBaseItem() GenericItem {
	if r.InventoryItemID.Valid {
		return &r.InventoryItem
	}

	return &r.BatchRecipe
}

func (r *RecipeItem) SetBaseItem(item GenericItem) {
	switch item.(type) {
	case *InventoryItem:
		r.InventoryItem = *item.(*InventoryItem)
	case *Recipe:
		r.BatchRecipe = *item.(*Recipe)
	}
}

func (r *RecipeItem) GetID() uuid.UUID {
	return r.ID
}
func (r *RecipeItem) GetBaseItemID() uuid.UUID {
	return r.GetBaseItem().GetID()
}

func (r *RecipeItem) GetName() string {
	return r.GetBaseItem().GetName()
}
func (r *RecipeItem) GetCategory() ItemCategory {
	return r.GetBaseItem().GetCategory()
}
func (r *RecipeItem) GetCategoryID() uuid.UUID {
	return r.GetBaseItem().GetCategoryID()
}
func (r *RecipeItem) SetCategory(category ItemCategory) {
	r.GetBaseItem().SetCategory(category)
}
func (r *RecipeItem) GetCountUnit() string {
	return r.GetBaseItem().GetCountUnit()
}
func (r *RecipeItem) GetIndex() int {
	return r.GetBaseItem().GetIndex()
}
func (r *RecipeItem) SetIndex(idx int) {
	r.Recipe.Index = idx
}
func (r *RecipeItem) GetRecipeUnit() string {
	if r.InventoryItemID.Valid {
		return r.InventoryItem.RecipeUnit
	}
	return r.BatchRecipe.RecipeUnit
}
func (r *RecipeItem) GetRecipeUnitConversion() float64 {
	if r.InventoryItemID.Valid {
		return r.InventoryItem.RecipeUnitConversion
	}
	return r.BatchRecipe.RecipeUnitConversion
}

func (r RecipeItem) GetSortValue() int {
	if r.InventoryItemID.Valid {
		return r.InventoryItem.GetSortValue()
	}

	return r.BatchRecipe.GetSortValue()
}

// Sort sorts the items based on category then inventory item indices
func (r *RecipeItems) Sort() {
	sort.Slice(*r, func(i, j int) bool {
		return (*r)[i].GetSortValue() < (*r)[j].GetSortValue()
	})
}

func (r *RecipeItems) ToGenericItems() *[]GenericItem {
	items := make([]GenericItem, len(*r))
	for i := 0; i < len(*r); i++ {
		items[i] = &(*r)[i]
	}

	return &items
}

func (r *RecipeItems) ToCompoundItems() *[]CompoundItem {
	items := make([]CompoundItem, len(*r))
	for i := 0; i < len(*r); i++ {
		items[i] = &(*r)[i]
	}

	return &items
}

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
		&validators.StringIsPresent{Field: r.Measure, Name: "Measure"},
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

func (r RecipeItem) getBaseItem() GenericItem {
	if r.InventoryItemID.Valid {
		return r.InventoryItem
	}

	return r.BatchRecipe
}

func (r RecipeItem) GetID() uuid.UUID {
	return r.ID
}
func (r RecipeItem) GetInventoryItemID() uuid.UUID {
	return r.getBaseItem().GetID()
}
func (r RecipeItem) GetName() string {
	return r.getBaseItem().GetName()
}
func (r RecipeItem) GetCategory() Category {
	return r.getBaseItem().GetCategory()
}
func (r RecipeItem) GetCountUnit() string {
	return r.getBaseItem().GetCountUnit()
}
func (r RecipeItem) GetIndex() int {
	return r.getBaseItem().GetIndex()
}

func (r RecipeItem) GetSortValue() int {
	if r.InventoryItemID.Valid {
		return r.InventoryItem.GetSortValue()
	}
	// HACK: bit of a hack - want recipes to always be after inv items
	return r.BatchRecipe.GetSortValue() + 100000
}

// Sort sorts the items based on category then inventory item indices
func (r *RecipeItems) Sort() {
	sort.Slice(*r, func(i, j int) bool {
		return (*r)[i].GetSortValue() < (*r)[j].GetSortValue()
	})
}

package models

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type InventoryItem struct {
	ID                   uuid.UUID    `json:"id" db:"id"`
	CreatedAt            time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at" db:"updated_at"`
	Name                 string       `json:"name" db:"name"`
	Category             ItemCategory `belongs_to:"item_categories" db:"-"`
	CategoryID           uuid.UUID    `json:"category_id" db:"category_id"`
	CountUnit            string       `json:"count_unit" db:"count_unit"`
	RecipeUnit           string       `json:"recipe_unit" db:"recipe_unit"`
	RecipeUnitConversion float64      `json:"recipe_unit_conversion" db:"recipe_unit_conversion"`
	Yield                float64      `json:"yield" db:"yield"`
	Index                int          `json:"index" db:"index"`
	IsActive             bool         `json:"is_active" db:"is_active"`
}

// String is not required by pop and may be deleted
func (i InventoryItem) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// InventoryItems is not required by pop and may be deleted
type InventoryItems []InventoryItem

// String is not required by pop and may be deleted
func (i InventoryItems) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *InventoryItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Name, Name: "Name"},
		&validators.StringIsPresent{Field: i.CountUnit, Name: "CountUnit"},
		&validators.StringIsPresent{Field: i.RecipeUnit, Name: "RecipeUnit"},
		&validators.IntIsGreaterThan{
			Field:    i.Index,
			Name:     "Index",
			Compared: -1,
		},
		&validators.FuncValidator{
			Field:   "",
			Name:    "Yield",
			Message: "Yield must be greater than 1",
			Fn: func() bool {
				return i.Yield > 0
			},
		},
		&validators.FuncValidator{
			Field:   "",
			Name:    "RecipeUnitConversion",
			Message: "Recipe Unit Conversion must be greater than 0",
			Fn: func() bool {
				return i.RecipeUnitConversion > 0
			},
		},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *InventoryItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return i.validateUniqueName(tx.Q())
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *InventoryItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	item := &InventoryItem{}
	if err := tx.Find(item, i.ID); err != nil {
		return validate.NewErrors(), err
	}
	if item.Name != i.Name {
		query := tx.Where("id != ?", i.ID)
		return i.validateUniqueName(query)
	}
	return validate.NewErrors(), nil
}

func (i *InventoryItem) validateUniqueName(query *pop.Query) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	items := &InventoryItems{}
	if err := query.All(items); err != nil {
		return verrs, err
	}

	for _, item := range *items {
		if strings.ToLower(i.Name) == strings.ToLower(item.Name) {
			verrs.Add("Name", "Item name already exists")
			break
		}
	}

	return verrs, nil
}

func (i *InventoryItem) GetID() uuid.UUID {
	return i.ID
}

func (i *InventoryItem) GetInventoryItemID() uuid.UUID {
	return i.ID
}

func (i *InventoryItem) GetBaseItem() GenericItem {
	return i
}

func (i *InventoryItem) GetName() string {
	return i.Name
}

func (i *InventoryItem) GetCategory() ItemCategory {
	return i.Category
}

func (i *InventoryItem) GetCountUnit() string {
	return i.CountUnit
}

func (i *InventoryItem) GetIndex() int {
	return i.Index
}

// GetSortValue returns a value for sorting where Category is highest prcedence
// and item index is second
func (i InventoryItem) GetSortValue() int {
	return i.Category.Index*1000 + i.Index
}

// Sort sorts the items based on category then inventory item indices
func (inv *InventoryItems) Sort() {
	sort.Slice(*inv, func(i, j int) bool {
		return (*inv)[i].GetSortValue() < (*inv)[j].GetSortValue()
	})
}

func (inv *InventoryItems) ToGenericItems() *[]GenericItem {
	items := make([]GenericItem, len(*inv))
	for i := 0; i < len(*inv); i++ {
		items[i] = &(*inv)[i]
	}

	return &items
}

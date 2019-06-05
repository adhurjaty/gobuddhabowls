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

type PrepItem struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	CountUnit     string    `json:"count_unit" db:"count_unit"`
	BatchRecipeID uuid.UUID `json:"batch_recipe_id" db:"batch_recipe_id"`
	BatchRecipe   Recipe    `belongs_to:"recipes" db:"-"`
	// Conversion is the number of recipe units in a prep item count
	Conversion float64 `json:"conversion" db:"conversion"`
	Index      int     `json:"index" db:"index"`
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

func (p *PrepItem) GetID() uuid.UUID {
	return p.ID
}

func (p *PrepItem) GetIndex() int {
	return p.Index
}

func (p *PrepItem) GetBaseItemID() uuid.UUID {
	return p.BatchRecipeID
}

func (p *PrepItem) GetBaseItem() GenericItem {
	return &p.BatchRecipe
}

func (p *PrepItem) SetBaseItem(item GenericItem) {
	p.BatchRecipe = *item.(*Recipe)
}

func (p *PrepItem) GetSortValue() int {
	return p.BatchRecipe.Category.Index*1000 + p.Index
}

func (p *PrepItem) GetName() string {
	return p.GetBaseItem().GetName()
}

func (p *PrepItem) GetCategory() ItemCategory {
	return p.GetBaseItem().GetCategory()
}

func (p *PrepItem) GetCountUnit() string {
	return p.CountUnit
}

func (p *PrepItems) ToModels() *[]Model {
	models := make([]Model, len(*p))
	for idx := range *p {
		models[idx] = &(*p)[idx]
	}

	return &models
}

func (p *PrepItems) Sort() {
	sort.Slice(*p, func(i, j int) bool {
		return (*p)[i].GetSortValue() < (*p)[j].GetSortValue()
	})
}

func (p *PrepItems) ToGenericItems() *[]GenericItem {
	items := make([]GenericItem, len(*p))
	for i := 0; i < len(*p); i++ {
		items[i] = &(*p)[i]
	}

	return &items
}

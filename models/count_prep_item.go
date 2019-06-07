package models

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type CountPrepItem struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	PrepItemID  uuid.UUID `json:"prep_item_id" db:"prep_item_id"`
	PrepItem    PrepItem  `belongs_to:"prep_items" db:"-"`
	InventoryID uuid.UUID `json:"inventory_id" db:"inventory_id"`
	Inventory   Inventory `belongs_to:"inventories" db:"-"`
	Count       float64   `json:"count" db:"count"`
}

// String is not required by pop and may be deleted
func (c CountPrepItem) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// CountPrepItems is not required by pop and may be deleted
type CountPrepItems []CountPrepItem

// String is not required by pop and may be deleted
func (c CountPrepItems) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *CountPrepItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *CountPrepItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *CountPrepItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (c *CountPrepItem) GetID() uuid.UUID {
	return c.ID
}

func (c *CountPrepItem) GetName() string {
	return c.GetBaseItem().GetName()
}

func (c *CountPrepItem) GetCategory() ItemCategory {
	return c.GetBaseItem().GetCategory()
}

func (c *CountPrepItem) SetCategory(category ItemCategory) {
	c.GetBaseItem().SetCategory(category)
}

func (c *CountPrepItem) GetCountUnit() string {
	return c.GetBaseItem().GetCountUnit()
}

func (c *CountPrepItem) GetIndex() int {
	return c.GetBaseItem().GetIndex()
}

func (c *CountPrepItem) GetBaseItemID() uuid.UUID {
	return c.PrepItemID
}

func (c *CountPrepItem) GetBaseItem() GenericItem {
	return &c.PrepItem
}

func (c *CountPrepItem) SetBaseItem(item GenericItem) {
	c.PrepItem = *item.(*PrepItem)
}

func (c *CountPrepItems) ToModels() []Model {
	items := make([]Model, len(*c))
	for i := range *c {
		items[i] = &(*c)[i]
	}

	return items
}

func (c *CountPrepItems) ToGenericItems() *[]GenericItem {
	items := make([]GenericItem, len(*c))
	for i := range *c {
		items[i] = &(*c)[i]
	}

	return &items
}

func (c *CountPrepItems) ToCompoundItems() *[]CompoundItem {
	items := make([]CompoundItem, len(*c))
	for i := range *c {
		items[i] = &(*c)[i]
	}

	return &items
}

func (c *CountPrepItems) Sort() {
	sort.Slice(*c, func(i, j int) bool {
		return (*c)[i].PrepItem.GetSortValue() < (*c)[j].PrepItem.GetSortValue()
	})
}

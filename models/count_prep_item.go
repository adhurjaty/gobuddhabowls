package models

import (
	"encoding/json"
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
	LineCount   float64   `json:"line_count" db:"line_count"`
	WalkInCount float64   `json:"walk_in_count" db:"walk_in_count"`
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

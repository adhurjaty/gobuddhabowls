package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type ItemCategory struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Name       string    `json:"name" db:"name"`
	Background string    `json:"background" db:"background"`
	Index      int       `json:"index" db:"index"`
}

// String is not required by pop and may be deleted
func (i ItemCategory) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// ItemCategories is not required by pop and may be deleted
type ItemCategories []ItemCategory

// String is not required by pop and may be deleted
func (i ItemCategories) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *ItemCategory) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Name, Name: "Name"},
		&validators.StringIsPresent{Field: i.Background, Name: "Background"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *ItemCategory) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	// TODO: unique name checker
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *ItemCategory) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (c ItemCategory) GetID() uuid.UUID {
	return c.ID
}

func (c ItemCategory) GetName() string {
	return c.Name
}

func (c ItemCategory) GetBackground() string {
	return c.Background
}

func (c ItemCategory) GetIndex() int {
	return c.Index
}

func (c ItemCategory) SetIndex(idx int) {
	c.Index = idx
}

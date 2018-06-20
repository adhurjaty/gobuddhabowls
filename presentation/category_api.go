package presentation

import (
	"buddhabowls/models"
	"errors"
)

// CategoryAPI category object for ui
type CategoryAPI struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Background string `json:"background"`
	Index      int    `json:"index"`
}

// ConvertToAPI converts a category to an api category
func (c *CategoryAPI) ConvertToAPI(m interface{}) error {
	category, ok := m.(models.InventoryItemCategory)
	if !ok {
		return errors.New("Must supply InventoryItemCategory type")
	}

	c.ID = category.ID.String()
	c.Name = category.Name
	c.Background = category.Background
	c.Index = category.Index

	return nil
}

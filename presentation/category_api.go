package presentation

import (
	"buddhabowls/models"
)

// CategoryAPI category object for ui
type CategoryAPI struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Background string `json:"background"`
	Index      int    `json:"index"`
}

// NewCategoryAPI converts a category to an api category
func NewCategoryAPI(category models.InventoryItemCategory) CategoryAPI {
	c := CategoryAPI{}

	c.ID = category.ID.String()
	c.Name = category.Name
	c.Background = category.Background
	c.Index = category.Index

	return c
}

package actions

import (
	"buddhabowls/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// InventoryItemCategoriesResource is a resource for managing InventoryItemCategories
type InventoryItemCategoriesResource struct {
	buffalo.Resource
}

// List serves the page for editing the categories for inventory items (color, order)
// GET /inventory_item_categories
func (v InventoryItemCategoriesResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	categories := &models.InventoryItemCategories{}
	if err := tx.Eager().All(categories); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.Auto(c, categories))
}

// Update updates the selected category
// PUT /inventory_item_categories/{purchase_order_id}
func (v InventoryItemCategoriesResource) Update(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

// Show gets the data for one category. Probably won't be used This function is mapped to
// the path GET /inventory_item_categories/{purchase_order_id}
func (v InventoryItemCategoriesResource) Show(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

func (v InventoryItemCategoriesResource) New(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

func (v InventoryItemCategoriesResource) Create(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

func (v InventoryItemCategoriesResource) Edit(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

func (v InventoryItemCategoriesResource) Destroy(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

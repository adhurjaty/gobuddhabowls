package actions

import (
	"buddhabowls/models"
	"buddhabowls/presentation"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
)

var _ = fmt.Println

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

	presenter := presentation.NewPresenter(tx)
	categories, err := presenter.GetAllCategories()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("inventoryItemCategories", categories)
	return c.Render(200, r.HTML("inventory_item_categories/index"))
}

// Update updates the selected category
// PUT /inventory_item_categories/{inventory_item_category_id}
func (v InventoryItemCategoriesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty InventoryItemCategory
	category := &models.InventoryItemCategory{}

	newItem := false

	if err := tx.Find(category, c.Param("inventory_item_category_id")); err != nil {
		newItem = true
	}

	// Bind InventoryItemCategory to the html form elements
	err := c.Bind(category)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("!!!!!!!!!!!!!!!!!!!")
	fmt.Println(category)
	fmt.Println(newItem)

	var verrs *validate.Errors
	if newItem {
		verrs, err = tx.ValidateAndCreate(category)
	} else {
		verrs, err = tx.ValidateAndUpdate(category)
	}
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Error(422, errors.New(
			fmt.Sprintf("Validation errors on categories: %v", verrs)))
	}

	// If there are no errors set a success message
	// c.Flash().Add("success", "PurchaseOrder was updated successfully")

	// and redirect to the purchase_orders index page
	return c.Render(200, r.String("success"))
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
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty InventoryItemCategory
	category := &models.InventoryItemCategory{}

	err := tx.Find(category, c.Param("inventory_item_category_id"))

	if err == nil {
		err = tx.Destroy(category)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	presenter := presentation.NewPresenter(tx)
	categories, err := presenter.GetAllCategories()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("inventoryItemCategories", categories)
	return c.Render(200, r.HTML("inventory_item_categories/index"))
}

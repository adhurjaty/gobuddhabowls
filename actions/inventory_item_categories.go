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

// ItemCategoriesResource is a resource for managing ItemCategories
type ItemCategoriesResource struct {
	buffalo.Resource
}

// List serves the page for editing the categories for inventory items (color, order)
// GET /inventory_item_categories
func (v ItemCategoriesResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	categories, err := presenter.GetAllCategories()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("itemCategories", categories)
	return c.Render(200, r.HTML("inventory_item_categories/index"))
}

// Update updates the selected category
// PUT /inventory_item_categories/{item_category_id}
func (v ItemCategoriesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty ItemCategory
	category := &models.ItemCategory{}

	newItem := false

	if err := tx.Find(category, c.Param("item_category_id")); err != nil {
		newItem = true
	}

	// Bind ItemCategory to the html form elements
	err := c.Bind(category)
	if err != nil {
		return errors.WithStack(err)
	}

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
func (v ItemCategoriesResource) Show(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

func (v ItemCategoriesResource) New(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

func (v ItemCategoriesResource) Create(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

func (v ItemCategoriesResource) Edit(c buffalo.Context) error {
	return c.Render(404, r.String("Not implemented"))
}

func (v ItemCategoriesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty ItemCategory
	category := &models.ItemCategory{}

	err := tx.Find(category, c.Param("item_category_id"))

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

	c.Set("itemCategories", categories)
	return c.Render(200, r.HTML("inventory_item_categories/index"))
}

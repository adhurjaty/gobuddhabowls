package actions

import (
	"buddhabowls/models"
	"buddhabowls/presentation"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

var _ = fmt.Printf

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (InventoryItem)
// DB Table: Plural (inventory_items)
// Resource: Plural (InventoryItems)
// Path: Plural (/inventory_items)
// View Template Folder: Plural (/templates/inventory_items/)

// InventoryItemsResource is the resource for the InventoryItem model
type InventoryItemsResource struct {
	buffalo.Resource
}

// List gets all InventoryItems. This function is mapped to the path
// GET /inventory_items
func (v InventoryItemsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	itemsAPI, err := presenter.GetInventoryItems()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("items", itemsAPI)

	return c.Render(200, r.HTML("inventory_items/index"))
}

// Show gets the data for one InventoryItem. This function is mapped to
// the path GET /inventory_items/{inventory_item_id}
func (v InventoryItemsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty InventoryItem
	inventoryItem := &models.InventoryItem{}

	// To find the InventoryItem the parameter inventory_item_id is used.
	if err := tx.Find(inventoryItem, c.Param("inventory_item_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, inventoryItem))
}

// New renders the form for creating a new InventoryItem.
// This function is mapped to the path GET /inventory_items/new
func (v InventoryItemsResource) New(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	vendorItems, err := presenter.GetBlankVendorItems()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("vendorItems", vendorItems)
	return c.Render(200, r.Auto(c, &models.InventoryItem{}))
}

// Create adds a InventoryItem to the DB. This function is mapped to the
// path POST /inventory_items
func (v InventoryItemsResource) Create(c buffalo.Context) error {
	// Allocate an empty InventoryItem
	inventoryItem := &models.InventoryItem{}

	// Bind inventoryItem to the html form elements
	if err := c.Bind(inventoryItem); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(inventoryItem)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, inventoryItem))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "InventoryItem was created successfully")

	// and redirect to the inventory_items index page
	return c.Render(201, r.Auto(c, inventoryItem))
}

// Edit renders a edit form for a InventoryItem. This function is
// mapped to the path GET /inventory_items/{inventory_item_id}/edit
func (v InventoryItemsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty InventoryItem
	inventoryItem := &models.InventoryItem{}

	if err := tx.Find(inventoryItem, c.Param("inventory_item_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, inventoryItem))
}

// Update changes a InventoryItem in the DB. This function is mapped to
// the path PUT /inventory_items/{inventory_item_id}
func (v InventoryItemsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	item, err := presenter.GetInventoryItem(c.Param("inventory_item_id"))
	if err != nil {
		return c.Error(404, err)
	}

	// Bind InventoryItem to the html form elements
	if err := c.Bind(item); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := presenter.UpdateInventoryItem(item)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// TODO: send error to page
		return c.Render(422, r.String("failure"))
	}

	// and redirect to the inventory_items index page
	return c.Render(200, r.String("success"))
}

// Destroy deletes a InventoryItem from the DB. This function is mapped
// to the path DELETE /inventory_items/{inventory_item_id}
func (v InventoryItemsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty InventoryItem
	inventoryItem := &models.InventoryItem{}

	// To find the InventoryItem the parameter inventory_item_id is used.
	if err := tx.Find(inventoryItem, c.Param("inventory_item_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(inventoryItem); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "InventoryItem was destroyed successfully")

	// Redirect to the inventory_items index page
	return c.Render(200, r.Auto(c, inventoryItem))
}

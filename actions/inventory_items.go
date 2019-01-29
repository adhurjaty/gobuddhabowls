package actions

import (
	"buddhabowls/helpers"
	"buddhabowls/models"
	"buddhabowls/presentation"
	"encoding/json"
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
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	items, err := presenter.GetMasterInventoryList()
	if err != nil {
		return err
	}

	inventory := &presentation.InventoryAPI{
		Date:  helpers.Today(),
		Items: *items,
	}

	c.Set("inventory", inventory)

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

	inventoryItem := &presentation.ItemAPI{
		Yield:                1,
		RecipeUnitConversion: 1,
	}
	presenter := presentation.NewPresenter(tx)
	err := setNewInvItemsViewVars(c, presenter)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("inventoryItem", inventoryItem)
	return c.Render(200, r.HTML("inventory_items/new"))
}

// Create adds a InventoryItem to the DB. This function is mapped to the
// path POST /inventory_items
func (v InventoryItemsResource) Create(c buffalo.Context) error {
	// Allocate an empty InventoryItem
	inventoryItem := &presentation.ItemAPI{}

	// Bind inventoryItem to the html form elements
	if err := bindInventoryItem(c, inventoryItem); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)

	// Validate the data from the html form
	verrs, err := presenter.InsertInventoryItem(inventoryItem)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		err = setNewInvItemsViewVars(c, presenter)
		if err != nil {
			return errors.WithStack(err)
		}
		c.Set("inventoryItem", inventoryItem)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(200, r.HTML("inventory_items/new"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "InventoryItem was created successfully")

	// and redirect to the inventory_items index page
	return c.Redirect(303, "/inventories")
}

func bindInventoryItem(c buffalo.Context, inventoryItem *presentation.ItemAPI) error {
	if err := c.Bind(inventoryItem); err != nil {
		return err
	}

	catID := c.Param("CategoryID")
	if catID != "" {
		inventoryItem.Category.ID = catID
	}

	vimStr := c.Param("VendorItemMap")
	if vimStr != "" {
		vendorItemsMap := make(map[string]presentation.ItemAPI)
		if err := json.Unmarshal([]byte(vimStr), &vendorItemsMap); err != nil {
			return err
		}
		inventoryItem.VendorItemMap = vendorItemsMap
	}

	return nil
}

func setNewInvItemsViewVars(c buffalo.Context, presenter *presentation.Presenter) error {
	vendorItems, err := presenter.GetBlankVendorItems()
	if err != nil {
		return err
	}
	categories, err := presenter.GetAllCategories()
	if err != nil {
		return err
	}
	inventoryItems, err := presenter.GetInventoryItems()
	if err != nil {
		return err
	}

	c.Set("vendorItems", vendorItems)
	c.Set("categories", categories)
	c.Set("inventoryItems", inventoryItems)

	return nil
}

// Edit renders a edit form for a InventoryItem. This function is
// mapped to the path GET /inventory_items/{inventory_item_id}/edit
func (v InventoryItemsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	inventoryItem, err := presenter.GetFullInventoryItem(c.Param("inventory_item_id"))
	if err != nil {
		return c.Error(404, err)
	}

	if err = setEditInvItemsViewVars(c, presenter, inventoryItem); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.HTML("inventory_items/edit"))
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
	if err := bindInventoryItem(c, item); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := presenter.UpdateFullInventoryItem(item)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		if err = setEditInvItemsViewVars(c, presenter, item); err != nil {
			return errors.WithStack(err)
		}

		// TODO: send error to page
		return c.Render(422, r.HTML("/inventories/edit"))
	}

	// and redirect to the inventory_items index page
	return c.Redirect(303, "/inventory_items")
	// return c.Render(200, r.String("SFEe"))
}

func setEditInvItemsViewVars(c buffalo.Context, presenter *presentation.Presenter,
	inventoryItem *presentation.ItemAPI) error {

	vendorItems, err := presenter.GetBlankVendorItems()
	if err != nil {
		return errors.WithStack(err)
	}
	for i, item := range *vendorItems {
		for vendorName, vendItem := range inventoryItem.VendorItemMap {
			if vendorName == item.SelectedVendor {
				(*vendorItems)[i] = vendItem
				break
			}
		}
	}

	categories, err := presenter.GetAllCategories()
	if err != nil {
		return errors.WithStack(err)
	}
	inventoryItems, err := presenter.GetInventoryItems()
	if err != nil {
		return errors.WithStack(err)
	}
	for i, item := range *inventoryItems {
		if item.ID == inventoryItem.ID {
			*inventoryItems = append((*inventoryItems)[:i],
				(*inventoryItems)[i+1:]...)
		}
	}

	c.Set("inventoryItems", inventoryItems)
	c.Set("categories", categories)
	c.Set("vendorItems", vendorItems)
	c.Set("inventoryItem", inventoryItem)

	return nil
}

func UpdateInventoryItem(c buffalo.Context) error {
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

	presenter := presentation.NewPresenter(tx)
	inventoryItem, err := presenter.GetFullInventoryItem(c.Param("inventory_item_id"))
	if err != nil {
		return c.Error(404, err)
	}

	if err = presenter.DestroyOrDeactivateInventoryItem(inventoryItem); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "InventoryItem was destroyed successfully")

	// Redirect to the inventory_items index page
	return c.Redirect(303, "/inventories")
}

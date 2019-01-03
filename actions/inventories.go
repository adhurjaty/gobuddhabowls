package actions

import (
	"buddhabowls/helpers"
	"buddhabowls/models"
	"buddhabowls/presentation"
	"fmt"
	"time"

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
// Model: Singular (Inventory)
// DB Table: Plural (inventories)
// Resource: Plural (Inventories)
// Path: Plural (/inventories)
// View Template Folder: Plural (/templates/inventories/)

// InventoriesResource is the resource for the Inventory model
type InventoriesResource struct {
	buffalo.Resource
}

// List gets master inventory
// This function is mapped to the path GET /inventories
func (v InventoriesResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	items, err := presenter.GetNewInventoryItems()
	if err != nil {
		return err
	}

	inventory := &presentation.InventoryAPI{
		Date:  helpers.Today(),
		Items: *items,
	}

	c.Set("inventory", inventory)

	return c.Render(200, r.HTML("inventories/index"))
}

// History gets all Inventories. This function is mapped to the path
// GET /inventories/history
func (v InventoriesResource) History(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	if err := setInventoryListVars(c, presenter); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.HTML("inventories/history"))
}

func setInventoryListVars(c buffalo.Context, presenter *presentation.Presenter) error {
	startTime, endTime, err := setPeriodSelector(c, presenter)
	if err != nil {
		return err
	}

	inventories, err := presenter.GetInventories(startTime, endTime)
	if err != nil {
		return err
	}

	if len(*inventories) == 0 {
		latestInv, err := presenter.GetLatestInventory(startTime)
		if err != nil {
			return err
		}

		*inventories = append(*inventories, *latestInv)
	}

	// fmt.Println(inventories)
	c.Set("inventories", inventories)
	c.Set("defaultInventory", (*inventories)[0])

	return nil
}

// Show gets the data for one Inventory. This function is mapped to
// the path GET /inventories/{inventory_id}
func (v InventoriesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Inventory
	inventory := &models.Inventory{}

	// To find the Inventory the parameter inventory_id is used.
	if err := tx.Find(inventory, c.Param("inventory_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, inventory))
}

// New renders the form for creating a new Inventory.
// This function is mapped to the path GET /inventories/new
func (v InventoriesResource) New(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	items, err := presenter.GetNewInventoryItems()
	if err != nil {
		return err
	}

	inventory := &presentation.InventoryAPI{
		Date:  helpers.Today(),
		Items: *items,
	}

	c.Set("inventory", inventory)

	return c.Render(200, r.HTML("inventories/new"))
}

// Create adds a Inventory to the DB. This function is mapped to the
// path POST /inventories
func (v InventoriesResource) Create(c buffalo.Context) error {
	// Allocate an empty Inventory
	invAPI := &presentation.InventoryAPI{}

	// Bind inventory to the html form elements
	err := c.Bind(invAPI)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	itemsParamJSON := c.Request().Form.Get("Items")
	if itemsParamJSON != "" {
		invAPI.Items, err = getItemsFromParams(itemsParamJSON)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	presenter := presentation.NewPresenter(tx)
	verrs, err := presenter.InsertInventory(invAPI)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		if err = setInventoryListVars(c, presenter); err != nil {
			return errors.WithStack(err)
		}

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("inventories/new"))
	}

	// If there are no errors set a success message

	week := presenter.GetSelectedWeek(invAPI.Date)
	startTime := week.StartTime.Format(time.RFC3339)
	redirectURL := c.Request().URL
	q := redirectURL.Query()
	q.Add("StartTime", startTime)
	redirectURL.RawQuery = q.Encode()

	// If there are no errors set a success message
	c.Flash().Add("success", "Inventory was created successfully")

	// and redirect to the inventory index page
	return c.Redirect(303, c.Request().URL.String(), redirectURL.String())
}

// Edit renders a edit form for a Inventory. This function is
// mapped to the path GET /inventories/{inventory_id}/edit
func (v InventoriesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Inventory
	inventory := &models.Inventory{}

	if err := tx.Find(inventory, c.Param("inventory_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, inventory))
}

// Update changes a Inventory in the DB. This function is mapped to
// the path PUT /inventories/{inventory_id}
func (v InventoriesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	invAPI, err := presenter.GetInventory(c.Param("inventory_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	// Bind Inventory to the html form elements
	if err := c.Bind(invAPI); err != nil {
		return errors.WithStack(err)
	}

	itemsParamJSON := c.Request().Form.Get("Items")
	if itemsParamJSON != "" {
		invAPI.Items, err = getItemsFromParams(itemsParamJSON)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	verrs, err := presenter.UpdateInventory(invAPI)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		if err = setInventoryListVars(c, presenter); err != nil {
			return errors.WithStack(err)
		}

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("inventories"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Inventory was updated successfully")

	// and redirect to the inventories index page
	return c.Redirect(303, "/inventories")
	// return nil
}

// Destroy deletes a Inventory from the DB. This function is mapped
// to the path DELETE /inventories/{inventory_id}
func (v InventoriesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	invAPI, err := presenter.GetInventory(c.Param("inventory_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	if err = presenter.DestroyInventory(invAPI); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Inventory was destroyed successfully")

	// Redirect to the inventories index page
	return c.Redirect(303, "/inventories")
}

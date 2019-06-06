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

// Following naming logic is implemented in Buffalo:
// Model: Singular (Vendor)
// DB Table: Plural (vendors)
// Resource: Plural (Vendors)
// Path: Plural (/vendors)
// View Template Folder: Plural (/templates/vendors/)

// VendorsResource is the resource for the Vendor model
type VendorsResource struct {
	buffalo.Resource
}

// List gets all Vendors. This function is mapped to the path
// GET /vendors
func (v VendorsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	vendors, err := presenter.GetVendors()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("vendors", vendors)

	return c.Render(200, r.HTML("vendors/index"))
}

// Show gets the data for one Vendor. This function is mapped to
// the path GET /vendors/{vendor_id}
func (v VendorsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Vendor
	vendor := &models.Vendor{}

	// To find the Vendor the parameter vendor_id is used.
	if err := tx.Find(vendor, c.Param("vendor_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, vendor))
}

// New renders the form for creating a new Vendor.
// This function is mapped to the path GET /vendors/new
func (v VendorsResource) New(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no trainsaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	vendor := presentation.VendorAPI{}

	if err := setVendorFormVars(c, presenter, &vendor); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.HTML("vendors/new"))
}

// Create adds a Vendor to the DB. This function is mapped to the
// path POST /vendors
func (v VendorsResource) Create(c buffalo.Context) error {
	// Allocate an empty Vendor
	vendorAPI := &presentation.VendorAPI{}

	// Bind vendor to the html form elements
	err := c.Bind(vendorAPI)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	vendorAPI.Items, err = getItemsFromParams(c.Request().Form.Get("Items"))
	if err != nil {
		return errors.WithStack(err)
	}

	presenter := presentation.NewPresenter(tx)
	verrs, err := presenter.InsertVendor(vendorAPI)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		setVendorFormVars(c, presenter, vendorAPI)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("vendors/new"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Vendor was created successfully")

	// and redirect to the vendors index page
	return c.Redirect(303, "/vendors")
}

// Edit renders a edit form for a Vendor. This function is
// mapped to the path GET /vendors/{vendor_id}/edit
func (v VendorsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	vendor, err := presenter.GetVendor(c.Param("vendor_id"))
	if err != nil {
		return c.Error(404, err)
	}

	if err = setVendorFormVars(c, presenter, vendor); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.HTML("vendors/edit"))
}

// Update changes a Vendor in the DB. This function is mapped to
// the path PUT /vendors/{vendor_id}
func (v VendorsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	vendor, err := presenter.GetVendor(c.Param("vendor_id"))
	if err != nil {
		return c.Error(404, err)
	}

	// Bind Vendor to the html form elements
	if err := c.Bind(vendor); err != nil {
		return errors.WithStack(err)
	}

	itemsParamJSON := c.Request().Form.Get("Items")
	if itemsParamJSON != "" {
		vendor.Items, err = getItemsFromParams(itemsParamJSON)
		if err != nil {
			return err
		}
	}

	verrs, err := presenter.UpdateVendor(vendor)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		setVendorFormVars(c, presenter, vendor)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("vendors/edit"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Vendor was updated successfully")

	// and redirect to the vendors index page
	return c.Redirect(303, "/vendors")
}

func UpdateVendorInline(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	vendor, err := presenter.GetVendor(c.Param("vendor_id"))
	if err != nil {
		return c.Error(404, err)
	}

	// Bind Vendor to the html form elements
	if err := c.Bind(vendor); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := presenter.UpdateVendorNoItems(vendor)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.String("failure"))
	}

	return c.Render(200, r.String("success"))
}

// Destroy deletes a Vendor from the DB. This function is mapped
// to the path DELETE /vendors/{vendor_id}
func (v VendorsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Vendor
	vendor := &models.Vendor{}

	// To find the Vendor the parameter vendor_id is used.
	if err := tx.Find(vendor, c.Param("vendor_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(vendor); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Vendor was destroyed successfully")

	// Redirect to the vendors index page
	return c.Render(200, r.Auto(c, vendor))
}

func setVendorFormVars(c buffalo.Context, presenter *presentation.Presenter, vendor *presentation.VendorAPI) error {
	inventoryItems, err := presenter.GetInventoryItems()
	// sort.Slice(*inventoryItems, func(i, j int) bool {
	// 	return (*inventoryItems)[i].Name < (*inventoryItems)[j].Name
	// })
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("vendor", vendor)
	c.Set("inventoryItems", inventoryItems)

	return nil
}

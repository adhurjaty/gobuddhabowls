package actions

// TODO: Change this to Presenter-ViewModel pattern
// each handler should get a presenter object and c.Set() all
// appropriate variables

import (
	"buddhabowls/logic"
	"buddhabowls/models"
	"buddhabowls/presentation"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

var _ = fmt.Printf

// Following naming logic is implemented in Buffalo:
// Model: Singular (PurchaseOrder)
// DB Table: Plural (purchase_orders)
// Resource: Plural (PurchaseOrders)
// Path: Plural (/purchase_orders)
// View Template Folder: Plural (/templates/purchase_orders/)

// PurchaseOrdersResource is the resource for the PurchaseOrder model
type PurchaseOrdersResource struct {
	buffalo.Resource
}

const (
	_poStartTimeKey = "poStartTime"
	_poEndTimeKey   = "poEndTime"
)

// List gets all PurchaseOrders. This function is mapped to the path
// GET /purchase_orders
// optional params: StartTime, [EndTime]
func (v PurchaseOrdersResource) List(c buffalo.Context) error {

	// get the parameters from URL
	paramsMap, ok := c.Params().(url.Values)
	if !ok {
		return c.Error(500, errors.New("Could not parse params"))
	}

	startVal, startTimeExists := paramsMap["StartTime"]
	endVal, endTimeExists := paramsMap["EndTime"]
	startTime := time.Time{}
	endTime := time.Time{}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	var err error
	if startTimeExists {
		startTime, err = time.Parse(time.RFC3339, startVal[0])
		if err != nil {
			return errors.WithStack(err)
		}
	}

	presenter := presentation.NewPresenter(tx)
	if !startTimeExists && !endTimeExists {
		startTime = time.Now()
	}

	periods := presenter.GetPeriods(startTime)
	weeks := presenter.GetWeeks(startTime)
	years := presenter.GetYears()

	if endTimeExists {
		endTime, err = time.Parse(time.RFC3339, endVal[0])
		if err != nil {
			return errors.WithStack(err)
		}
		periods = append([]logic.Period{logic.Period{}}, periods...)
		weeks = append([]logic.Week{logic.Week{}}, weeks...)
		years = append([]int{0}, years...)

		c.Set("SelectedPeriod", periods[0])
		c.Set("SelectedWeek", weeks[0])
		c.Set("SelectedYear", startTime.Year())
	} else {
		selectedPeriod := presenter.GetSelectedPeriod(startTime)
		selectedWeek := presenter.GetSelectedWeek(startTime)
		c.Set("SelectedPeriod", selectedPeriod)
		c.Set("SelectedWeek", selectedWeek)
		c.Set("SelectedYear", startTime.Year())
		startTime = selectedWeek.StartTime
		endTime = selectedWeek.EndTime
	}

	purchaseOrders, err := presenter.GetPurchaseOrders(startTime, endTime)
	if err != nil {
		return errors.WithStack(err)
	}

	openPos := presentation.PurchaseOrdersAPI{}
	recPos := presentation.PurchaseOrdersAPI{}
	for _, po := range *purchaseOrders {
		if po.ReceivedDate.Valid {
			recPos = append(recPos, po)
		} else {
			openPos = append(openPos, po)
		}
	}

	c.Set("Periods", periods)
	c.Set("Weeks", weeks)
	c.Set("Years", years)
	c.Set("StartTime", startTime)
	c.Set("EndTime", endTime)
	c.Set("purchaseOrders", purchaseOrders)
	c.Set("openPurchaseOrders", openPos)
	c.Set("recPurchaseOrders", recPos)

	return c.Render(200, r.HTML("purchase_orders/index"))
}

// Show gets the data for one PurchaseOrder. This function is mapped to
// the path GET /purchase_orders/{purchase_order_id}
func (v PurchaseOrdersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty PurchaseOrder
	purchaseOrder := &models.PurchaseOrder{}

	// To find the PurchaseOrder the parameter purchase_order_id is used.
	if err := tx.Find(purchaseOrder, c.Param("purchase_order_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, purchaseOrder))
}

// New renders the form for creating a new PurchaseOrder.
// This function is mapped to the path GET /purchase_orders/new
func (v PurchaseOrdersResource) New(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	purchaseOrder := presentation.PurchaseOrderAPI{}

	err := setPurchaseOrderViewVars(c, presenter, purchaseOrder)
	if err != nil {
		return c.Error(500, err)
	}

	return c.Render(200, r.Auto(c, &models.PurchaseOrder{}))
}

// Create adds a PurchaseOrder to the DB. This function is mapped to the
// path POST /purchase_orders
func (v PurchaseOrdersResource) Create(c buffalo.Context) error {

	// Allocate an empty PurchaseOrder API object
	poAPI := &presentation.PurchaseOrderAPI{}

	// Bind purchaseOrder to the html form elements
	if err := c.Bind(poAPI); err != nil {
		return errors.WithStack(err)
	}
	vendorID := c.Request().Form.Get("VendorID")

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	vendor, err := presenter.GetVendor(vendorID)
	if err != nil {
		return errors.WithStack(err)
	}
	poAPI.Vendor = *vendor
	poAPI.ShippingCost = vendor.ShippingCost

	itemsJSON := c.Request().Form.Get("Items")
	poAPI.Items, err = getItemsFromParams(itemsJSON)
	if err != nil {
		return errors.WithStack(err)
	}

	verrs, err := presenter.InsertPurchaseOrder(poAPI)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		err = setPurchaseOrderViewVars(c, presenter, presentation.PurchaseOrderAPI{})
		if err != nil {
			return errors.WithStack(err)
		}
		c.Set("errors", verrs)
		return c.Render(422, r.Auto(c, models.PurchaseOrder{}))
	}

	week := presenter.GetSelectedWeek(poAPI.OrderDate.Time)
	startTime := week.StartTime.Format(time.RFC3339)
	redirectURL := c.Request().URL
	q := redirectURL.Query()
	q.Add("StartTime", startTime)
	redirectURL.RawQuery = q.Encode()

	// If there are no errors set a success message
	c.Flash().Add("success", "PurchaseOrder was created successfully")

	// and redirect to the purchase_orders index page
	return c.Redirect(303, c.Request().URL.String(), redirectURL.String())
}

// Edit renders a edit form for a PurchaseOrder. This function is
// mapped to the path GET /purchase_orders/{purchase_order_id}/edit
func (v PurchaseOrdersResource) Edit(c buffalo.Context) error {

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	purchaseOrder, err := presenter.GetPurchaseOrder(c.Param("purchase_order_id"))
	if err != nil {
		return c.Error(404, err)
	}

	err = setPurchaseOrderViewVars(c, presenter, *purchaseOrder)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.Auto(c, models.PurchaseOrder{}))
}

// Update changes a PurchaseOrder in the DB. This function is mapped to
// the path PUT /purchase_orders/{purchase_order_id}
func (v PurchaseOrdersResource) Update(c buffalo.Context) error {

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	poAPI, err := presenter.GetPurchaseOrder(c.Param("purchase_order_id"))
	if err != nil {
		return c.Error(404, err)
	}

	// Bind purchaseOrder to the html form elements
	if err := c.Bind(poAPI); err != nil {
		return errors.WithStack(err)
	}

	// invalidate ReceivedOrder if it is blank (default time value)
	if poAPI.ReceivedDate.Time.Unix() <= 0 ||
		c.Request().Form.Get("ReceivedDate") == "" {
		poAPI.ReceivedDate.Time = poAPI.OrderDate.Time
		poAPI.ReceivedDate.Valid = false
	}

	itemsParamJSON := c.Request().Form.Get("Items")
	if itemsParamJSON != "" {
		poAPI.Items, err = getItemsFromParams(itemsParamJSON)
		if err != nil {
			return err
		}
	}

	verrs, err := presenter.UpdatePurchaseOrder(poAPI)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		setPurchaseOrderViewVars(c, presenter, *poAPI)
		c.Set("errors", verrs)

		return c.Render(422, r.Auto(c, models.PurchaseOrder{}))
	}

	// need to check whether this is the most recent order from this vendor
	// TODO: move this to when order is received
	// selectedVendor, err := models.LoadVendor(tx, purchaseOrder.VendorID.String())
	// for _, vendorItem := range selectedVendor.Items {
	// 	if vendorItem.InventoryItemID == item.InventoryItemID {
	// 		if vendorItem.Price != item.Price { // && this is the most recent order from them
	// 			vendorItem.Price = item.Price
	// 			tx.ValidateAndUpdate(vendorItem)
	// 		}
	// 		break
	// 	}
	// }

	// // If there are no errors set a success message
	c.Flash().Add("success", "PurchaseOrder was updated successfully")

	week := presenter.GetSelectedWeek(poAPI.OrderDate.Time)
	startTime := week.StartTime.Format(time.RFC3339)
	redirectURL, _ := url.Parse("/purchase_orders")
	q := redirectURL.Query()
	q.Add("StartTime", startTime)
	redirectURL.RawQuery = q.Encode()

	// and redirect to the purchase_orders index page
	return c.Redirect(303, redirectURL.String())
}

// Destroy deletes a PurchaseOrder from the DB. This function is mapped
// to the path DELETE /purchase_orders/{purchase_order_id}
func (v PurchaseOrdersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	purchaseOrder, err := presenter.GetPurchaseOrder(c.Param("purchase_order_id"))
	if err != nil {
		return c.Error(404, err)
	}

	err = presenter.DestroyPurchaseOrder(purchaseOrder)
	if err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "PurchaseOrder was destroyed successfully")

	week := presenter.GetSelectedWeek(purchaseOrder.OrderDate.Time)
	startTime := week.StartTime.Format(time.RFC3339)
	redirectURL, _ := url.Parse("/purchase_orders")
	q := redirectURL.Query()
	q.Add("StartTime", startTime)
	redirectURL.RawQuery = q.Encode()

	// and redirect to the purchase_orders index page
	return c.Redirect(303, redirectURL.String())
}

// ShowOrderSheet displays the order sheet for sending to the vendor
// mapped to the path /purchase_orders/{purchase_order_id}/order_sheet
func ShowOrderSheet(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	purchaseOrder, err := presenter.GetPurchaseOrder(c.Param("purchase_order_id"))
	if err != nil {
		return c.Error(404, err)
	}
	vendor, err := presenter.GetVendor(purchaseOrder.Vendor.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	purchaseOrder.Items = presentation.AddVendorInfo(purchaseOrder.Items, vendor.Items)

	c.Set("purchaseOrder", purchaseOrder)

	return c.Render(200, r.HTML("purchase_orders/order_sheet.html", "old_application.html"))
}

// ShowReceivingList displays the receiving list for marking received orders at the restaurant
// mapped to the path /purchase_orders/{purchase_order_id}/receiving_list
func ShowReceivingList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	purchaseOrder, err := presenter.GetPurchaseOrder(c.Param("purchase_order_id"))
	if err != nil {
		return c.Error(404, err)
	}
	vendor, err := presenter.GetVendor(purchaseOrder.Vendor.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	purchaseOrder.Items = presentation.AddVendorInfo(purchaseOrder.Items, vendor.Items)

	c.Set("purchaseOrder", purchaseOrder)

	return c.Render(200, r.HTML("purchase_orders/receiving_list.html", "old_application.html"))
}

func setPurchaseOrderViewVars(c buffalo.Context, presenter *presentation.Presenter, poAPI presentation.PurchaseOrderAPI) error {
	newItem := poAPI.ID == ""

	var vendors *presentation.VendorsAPI
	if newItem {
		var err error
		vendors, err = presenter.GetVendors()
		if err != nil {
			return err
		}
		// add a blank vendor at the beginning so user is prompted to select a vendor
		vendorList := append(presentation.VendorsAPI{presentation.VendorAPI{}}, *vendors...)
		vendors = &vendorList
	} else {
		vendor, err := presenter.GetVendor(poAPI.Vendor.ID)
		if err != nil {
			return c.Error(500, err)
		}
		poAPI.Vendor = *vendor

		vendors = &presentation.VendorsAPI{poAPI.Vendor}

		remainingItems := presentation.ItemsAPI{}
		for _, vendorItem := range poAPI.Vendor.Items {
			contains := false
			for _, poItem := range poAPI.Items {
				if vendorItem.InventoryItemID == poItem.InventoryItemID {
					contains = true
					break
				}
			}
			if !contains {
				remainingItems = append(remainingItems, vendorItem)
			}
		}
		c.Set("remainingItems", remainingItems)
	}
	// map from vendor ID to vendor items
	vendorItemsMap := map[string]presentation.ItemsAPI{}
	for _, vendor := range *vendors {
		vendorItemsMap[vendor.ID] = vendor.Items
	}

	c.Set("po", poAPI)
	c.Set("vendors", vendors)
	c.Set("vendorItemsMap", vendorItemsMap)
	return nil
}

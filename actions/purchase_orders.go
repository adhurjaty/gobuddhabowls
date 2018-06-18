package actions

// TODO: Change this to Presenter-ViewModel pattern
// each handler should get a presenter object and c.Set() all
// appropriate variables

import (
	"buddhabowls/models"
	"buddhabowls/presentation"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

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

// PurchaseOrderDateChanged updates visible purchase orders table
// GET /purchase_orders/date_changed
// TODO: remove this and put into List
func PurchaseOrderDateChanged(c buffalo.Context) error {
	// store := c.Session()

	// // indicates whether user used the custome date range
	// customDateRange := false

	// // Get the DB connection from the context
	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return c.Error(500, errors.New("no transaction found"))
	// }

	// // get the parameters from URL
	// paramsMap, ok := c.Params().(url.Values)
	// if !ok {
	// 	return c.Error(500, errors.New("Could not parse params"))
	// }

	// yearStr, found := paramsMap["Year"]
	// if !found {
	// 	return c.Error(422, errors.New("No year supplied"))
	// }

	// year, err := strconv.Atoi(yearStr[0])
	// if err != nil {
	// 	return c.Error(500, errors.New("Could not get year from params"))
	// }
	// // if user selects the year option, then go to 1/1/year
	// startTime := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// endTime := startTime

	// // if user selected period, week or daterange
	// if startTimeStr, ok := paramsMap["StartTime"]; ok {
	// 	startTime, err = time.Parse(time.RFC3339, startTimeStr[0])
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf("Could not parse start time: %s", startTimeStr[0]))
	// 		return c.Error(500, errors.New("Could not parse start time"))
	// 	}

	// 	// if user selected daterange
	// 	if endTimeStr, ok := paramsMap["EndTime"]; ok {
	// 		endTime, err = time.Parse(time.RFC3339, endTimeStr[0])
	// 		if err != nil {
	// 			fmt.Println(fmt.Errorf("Could not parse end time: %s", endTimeStr[0]))
	// 			return c.Error(500, errors.New("Could not parse end time"))
	// 		}

	// 		customDateRange = true
	// 	}
	// }

	// periodSelectorContext := componentcontexts.PeriodSelectorContext{}
	// // TODO: make a new period selector function
	// periodSelectorContext.Init(startTime)

	// // only get the week if the date range was not selected
	// if !customDateRange {
	// 	startTime = periodSelectorContext.SelectedWeek.StartTime
	// 	endTime = periodSelectorContext.SelectedWeek.EndTime
	// }

	// startVal := startTime.Format(time.RFC3339)
	// endVal := endTime.Format(time.RFC3339)

	// // change selected PO dates in session
	// store.Set(_poStartTimeKey, startVal)
	// store.Set(_poEndTimeKey, endVal)

	// q := tx.Eager().Where("order_date >= ? AND order_date < ?",
	// 	startVal, endVal).Order("order_date DESC")

	// openPos, recPos, err := presentation.GetOpenRecPurchaseOrders(q)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	// // period selector view information
	// c.Set("pSelectorContext", periodSelectorContext)
	// c.Set("startTime", nulls.Time{Valid: true, Time: startTime})
	// c.Set("endTime", nulls.Time{Valid: true, Time: endTime})
	// c.Set("customDateRange", customDateRange)

	// // purchase order view information
	// c.Set("openPurchaseOrders", openPos)
	// c.Set("recPurchaseOrders", recPos)

	// lineChartData := presentation.GetLineChartJSONData(openPos, recPos)
	// c.Set("trendChartData", lineChartData)

	// // summary table view information
	// barChartData := presentation.GetBarChartJSONData(openPos, recPos)
	// c.Set("barChartData", barChartData)

	// years, ok := store.Get("years").([]int)
	// if !ok {
	// 	years, _ = models.GetYears(tx)
	// }
	// c.Set("years", years)

	return c.Render(200, r.JavaScript("purchase_orders/replace_table"))

}

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

	// indicates whether user used the custome date range
	customDateRange := !endTimeExists

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}
	presenter := presentation.NewPresenter(tx)
	if !startTimeExists && !endTimeExists {
		startTime = time.Now()
	}

	var err error
	if startTimeExists {
		unixTime, err := strconv.ParseInt(startVal[0], 10, 64)
		if err != nil {
			return errors.WithStack(err)
		}
		startTime = time.Unix(unixTime, 0)
	}
	if endTimeExists {
		unixTime, err := strconv.ParseInt(endVal[0], 10, 64)
		if err != nil {
			return errors.WithStack(err)
		}
		endTime = time.Unix(unixTime, 0)
	}

	periodSelector := presenter.GetPeriodContext(startTime)
	c.Set("pSelectorContext", periodSelector)
	c.Set("customDateRange", customDateRange)

	startTime = periodSelector.SelectedWeek.StartTime
	endTime = periodSelector.SelectedWeek.EndTime
	purchaseOrders, err := presenter.GetPurchaseOrders(startTime, endTime)
	if err != nil {
		return errors.WithStack(err)
	}
	c.Set("purchaseOrders", purchaseOrders)

	return c.Render(200, r.HTML("purchase_orders/index"))

	// EDIT POINT

	// indicates whether user used the custome date range
	// customDateRange := false

	// // Get the DB connection from the context
	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return errors.WithStack(errors.New("no transaction found"))
	// }

	// years, err := models.GetYears(tx)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	// store := c.Session()

	// periodSelectorContext := componentcontexts.PeriodSelectorContext{}
	// startTime := time.Time{}
	// endTime := startTime
	// if store.Get(_poStartTimeKey) == nil || store.Get(_poEndTimeKey) == nil {
	// 	periodSelectorContext.Init(time.Now())
	// 	startTime = periodSelectorContext.SelectedWeek.StartTime
	// 	endTime = periodSelectorContext.SelectedWeek.EndTime
	// 	store.Set(_poStartTimeKey, startTime.Format(time.RFC3339))
	// 	store.Set(_poEndTimeKey, endTime.Format(time.RFC3339))
	// } else {
	// 	t, err := time.Parse(time.RFC3339, store.Get(_poStartTimeKey).(string))
	// 	if err != nil {
	// 		return errors.WithStack(err)
	// 	}
	// 	startTime = t
	// 	endTime, err = time.Parse(time.RFC3339, store.Get(_poEndTimeKey).(string))
	// 	if err != nil {
	// 		return errors.WithStack(err)
	// 	}
	// 	periodSelectorContext.Init(startTime)

	// 	customDateRange = startTime != periodSelectorContext.SelectedWeek.StartTime ||
	// 		endTime != periodSelectorContext.SelectedWeek.EndTime

	// }

	// presenter := presentation.Presenter{}
	// periodData, err := presenter.GetPeriodData(tx)
	// if err != nil {
	// 	return err
	// }
	// c.Set("periodData", periodData)
	// c.Set("pSelectorContext", periodSelectorContext)
	// c.Set("startTime", nulls.Time{Valid: true, Time: startTime})
	// c.Set("endTime", nulls.Time{Valid: true, Time: endTime})
	// c.Set("years", years)
	// c.Set("customDateRange", customDateRange)

	// return c.Render(200, r.HTML("purchase_orders/index"))
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

	vendors := getSortedVendors(tx)

	c.Set("vendors", vendors)
	c.Set("purchaseOrder", models.PurchaseOrder{})

	return c.Render(200, r.Auto(c, &models.PurchaseOrder{}))
}

// NewOrderVendorChanged updates the new PO page when vendor has been selected
// GET /purchase_orders/order_vendor_changed/{vendor_id}
func NewOrderVendorChanged(c buffalo.Context) error {
	// Get the DB connection from the context
	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return errors.WithStack(errors.New("no transaction found"))
	// }

	// // get the vendor from params
	// selectedVendor, err := models.LoadVendor(tx, c.Param("vendor_id"))
	// if err != nil {
	// 	return c.Error(404, err)
	// }

	// // format vendor items to be shown in the UI
	// vendOrderItems := selectedVendor.Items.ToOrderItems()
	// categoryGroups := models.GetCategoryGroups(vendOrderItems.ToCountItems())

	// // get and sort keys from the map
	// sortedCategories := models.InventoryItemCategories{}
	// for k := range categoryGroups {
	// 	sortedCategories = append(sortedCategories, k)
	// }
	// sort.Slice(sortedCategories, func(i, j int) bool {
	// 	return sortedCategories[i].Index < sortedCategories[j].Index
	// })

	// // pass variables to UI
	// c.Set("sortedCategories", sortedCategories)
	// c.Set("categoryGroups", categoryGroups)

	return c.Render(200, r.JavaScript("purchase_orders/replace_new_vendor_items"))
}

// Create adds a PurchaseOrder to the DB. This function is mapped to the
// path POST /purchase_orders
func (v PurchaseOrdersResource) Create(c buffalo.Context) error {

	// Allocate an empty PurchaseOrder
	purchaseOrder := &models.PurchaseOrder{}

	// Bind purchaseOrder to the html form elements
	if err := c.Bind(purchaseOrder); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// get shipping cost from vendor
	vendor := &models.Vendor{}

	if err := tx.Find(vendor, purchaseOrder.VendorID); err != nil {
		return errors.WithStack(err)
	}
	purchaseOrder.ShippingCost = vendor.ShippingCost

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(purchaseOrder)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		vendors := getSortedVendors(tx)

		// Make the errors available inside the html template
		c.Set("errors", verrs)
		c.Set("vendors", vendors)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, purchaseOrder))
	}

	// Create the OrderItems as well
	itemsParamJSON, ok := c.Request().Form["Items"]
	if !ok {
		return c.Error(500, errors.New("Could not get items from form params"))
	}

	err = setItemsFromParams(itemsParamJSON[0], purchaseOrder)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, item := range purchaseOrder.Items {

		verrs, err := tx.ValidateAndCreate(&item)
		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			vendors := getSortedVendors(tx)

			// Make the errors available inside the html template
			c.Set("errors", verrs)
			c.Set("vendors", vendors)

			// Render again the new.html template that the user can
			// correct the input.
			return c.Render(422, r.Auto(c, purchaseOrder))
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
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "PurchaseOrder was created successfully")

	// and redirect to the purchase_orders index page
	return c.Render(201, r.Auto(c, purchaseOrder))
}

// PurchaseOrdersCountChanged updates the UI for new and edit orders when the count
// or price change. It displays a category price breakdown and updates
// the price extension
// mapped to the path POST /purchase_orders/count_changed
func PurchaseOrdersCountChanged(c buffalo.Context) error {
	purchaseOrder := &models.PurchaseOrder{}
	itemsJSON, ok := c.Request().Form["Items"]
	if !ok {
		return c.Error(500, errors.New("Could not get items from form"))
	}
	err := setItemsFromParams(itemsJSON[0], purchaseOrder)
	if err != nil {
		return errors.WithStack(err)
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return c.Error(500, errors.New("Could not open database connection"))
	}

	for i, item := range purchaseOrder.Items {
		invItem := &models.InventoryItem{}
		err = tx.Eager().Find(invItem, item.InventoryItemID)
		if err != nil {
			return c.Error(500, errors.New("invalid inventory item ID"))
		}
		purchaseOrder.Items[i].InventoryItem = *invItem
	}

	categoryDetails := purchaseOrder.GetCategoryCosts()
	c.Set("categoryDetails", categoryDetails)
	c.Set("title", "Category Breakdown")

	return c.Render(200, r.JavaScript("purchase_orders/count_changed"))
}

// Edit renders a edit form for a PurchaseOrder. This function is
// mapped to the path GET /purchase_orders/{purchase_order_id}/edit
func (v PurchaseOrdersResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return errors.WithStack(errors.New("no transaction found"))
	// }

	// purchaseOrder, err := models.LoadPurchaseOrder(tx, c.Param("purchase_order_id"))
	// if err != nil {
	// 	return c.Error(404, err)
	// }

	// setEditPOView(c, &purchaseOrder)

	// return c.Render(200, r.Auto(c, purchaseOrder))
	return c.Render(200, r.Auto(c, nil))
}

// AddPurchaseOrderItem adds an order item to the order item list
// mapped to /purchase_orders/add_item/{purchase_order_id}
func AddPurchaseOrderItem(c buffalo.Context) error {

	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return errors.WithStack(errors.New("no transaction found"))
	// }

	// vendorItemID := c.Request().Form["VendorItemID"][0]

	// vendorItem, err := models.LoadVendorItem(tx, vendorItemID)
	// if err != nil {
	// 	return c.Error(500, err)
	// }

	// purchaseOrder, err := models.LoadPurchaseOrder(tx, c.Param("purchase_order_id"))
	// if err != nil {
	// 	return c.Error(404, err)
	// }

	// // Bind PurchaseOrder to the html form elements
	// if err := c.Bind(&purchaseOrder); err != nil {
	// 	fmt.Println(err)
	// 	return errors.WithStack(err)
	// }

	// purchaseOrder.Items = append(purchaseOrder.Items, *vendorItem.ToOrderItem())
	// purchaseOrder.Items.Sort()

	// setEditPOView(c, &purchaseOrder)

	return c.Render(200, r.HTML("purchase_orders/edit"))
}

// Update changes a PurchaseOrder in the DB. This function is mapped to
// the path PUT /purchase_orders/{purchase_order_id}
func (v PurchaseOrdersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return errors.WithStack(errors.New("no transaction found"))
	// }

	// // Allocate an empty PurchaseOrder
	// purchaseOrder := &models.PurchaseOrder{}

	// if err := tx.Find(purchaseOrder, c.Param("purchase_order_id")); err != nil {
	// 	return c.Error(404, err)
	// }

	// // Bind PurchaseOrder to the html form elements
	// if err := c.Bind(purchaseOrder); err != nil {
	// 	fmt.Println(err)
	// 	return errors.WithStack(err)
	// }

	// verrs, err := tx.ValidateAndUpdate(purchaseOrder)
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Invalid data"))
	// 	return c.Error(422, err)
	// }
	// if verrs.HasAny() {
	// 	fmt.Println(fmt.Errorf("Invalid data"))
	// 	errorMsgs := []string{}
	// 	for _, verr := range verrs.Errors {
	// 		for _, v := range verr {
	// 			errorMsgs = append(errorMsgs, v)
	// 		}
	// 	}

	// 	return c.Render(422, r.String(strings.Join(errorMsgs, "\n")))
	// }

	// // Create the OrderItems as well
	// itemsParamJSON, ok := c.Request().Form["Items"]
	// if !ok {
	// 	// case for editing item within main datagrid
	// 	return c.Render(200, r.String("success"))
	// }

	// err = setItemsFromParams(itemsParamJSON[0], purchaseOrder)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	// for _, item := range purchaseOrder.Items {

	// 	verrs, err := tx.ValidateAndUpdate(&item)
	// 	if err != nil {
	// 		return errors.WithStack(err)
	// 	}

	// 	if verrs.HasAny() {
	// 		// models.LoadOrderItems(tx, purchaseOrder)
	// 		po, err := models.LoadPurchaseOrder(tx, purchaseOrder.ID.String())
	// 		if err != nil {
	// 			return errors.WithStack(err)
	// 		}
	// 		setEditPOView(c, &po)

	// 		// Render again the edit.html template that the user can
	// 		// correct the input.
	// 		return c.Render(422, r.Auto(c, &po))
	// 	}
	// 	// need to check whether this is the most recent order from this vendor
	// 	// TODO: move this to when order is received
	// 	// selectedVendor, err := models.LoadVendor(tx, purchaseOrder.VendorID.String())
	// 	// for _, vendorItem := range selectedVendor.Items {
	// 	// 	if vendorItem.InventoryItemID == item.InventoryItemID {
	// 	// 		if vendorItem.Price != item.Price { // && this is the most recent order from them
	// 	// 			vendorItem.Price = item.Price
	// 	// 			tx.ValidateAndUpdate(vendorItem)
	// 	// 		}
	// 	// 		break
	// 	// 	}
	// 	// }
	// }

	// // If there are no errors set a success message
	// c.Flash().Add("success", "PurchaseOrder was updated successfully")

	// and redirect to the purchase_orders index page
	// return c.Render(200, r.Auto(c, purchaseOrder))
	return c.Render(200, r.Auto(c, nil))
}

// Destroy deletes a PurchaseOrder from the DB. This function is mapped
// to the path DELETE /purchase_orders/{purchase_order_id}
func (v PurchaseOrdersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return errors.WithStack(errors.New("no transaction found"))
	// }

	// // To find the PurchaseOrder the parameter purchase_order_id is used.
	// purchaseOrder, err := models.LoadPurchaseOrder(tx, c.Param("purchase_order_id"))

	// if err != nil {
	// 	return c.Error(404, err)
	// }

	// // destroy associated order items
	// for _, item := range purchaseOrder.Items {
	// 	if err = tx.Destroy(&item); err != nil {
	// 		return errors.WithStack(err)
	// 	}
	// }

	// if err := tx.Destroy(&purchaseOrder); err != nil {
	// 	return errors.WithStack(err)
	// }

	// // If there are no errors set a flash message
	// c.Flash().Add("success", "PurchaseOrder was destroyed successfully")

	// // Redirect to the purchase_orders index page
	// return c.Render(200, r.Auto(c, purchaseOrder))
	return c.Render(200, r.Auto(c, nil))
}

func getSortedVendors(tx *pop.Connection) models.Vendors {
	vendors := models.Vendors{}
	tx.Eager().All(&vendors)

	// sort and add empty option
	sort.Slice(vendors, func(i, j int) bool {
		return vendors[i].Name < vendors[j].Name
	})
	vendors = append(models.Vendors{models.Vendor{}}, vendors...)

	return vendors
}

func setEditPOView(c buffalo.Context, purchaseOrder *models.PurchaseOrder) {
	c.Set("vendors", models.Vendors{purchaseOrder.Vendor})
	c.Set("purchaseOrder", purchaseOrder)

	categoryDetails := purchaseOrder.GetCategoryCosts()
	c.Set("categoryDetails", categoryDetails)
	c.Set("title", "Category Breakdown")

	categoryGroups := models.GetCategoryGroups(purchaseOrder.Items.ToCountItems())
	// get and sort keys from the map
	sortedCategories := models.InventoryItemCategories{}
	for k := range categoryGroups {
		sortedCategories = append(sortedCategories, k)
	}
	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Index < sortedCategories[j].Index
	})
	c.Set("sortedCategories", sortedCategories)
	c.Set("categoryGroups", categoryGroups)

	tx := c.Value("tx").(*pop.Connection)
	c.Set("remainingVendorItems", *getRemainingVendorItems(purchaseOrder, tx))
}

func setItemsFromParams(itemsParamJSON string, purchaseOrder *models.PurchaseOrder) error {
	orderItems := models.OrderItems{}

	err := json.Unmarshal([]byte(itemsParamJSON), &orderItems)
	if err != nil {
		return err
	}

	purchaseOrder.Items = models.OrderItems{}
	for _, item := range orderItems {
		if item.Count > 0 {
			item.OrderID = purchaseOrder.ID
			purchaseOrder.Items = append(purchaseOrder.Items, item)
		}
	}

	return nil
}

func getRemainingVendorItems(po *models.PurchaseOrder, tx *pop.Connection) *models.VendorItems {
	// vendorItems := models.VendorItems{}
	// vendor, err := models.LoadVendor(tx, po.VendorID.String())
	// if err != nil {
	// 	return nil
	// }

	// for _, vendorItem := range vendor.Items {
	// 	contains := func(vendorItem models.VendorItem) bool {
	// 		for _, orderItem := range po.Items {
	// 			if orderItem.InventoryItemID == vendorItem.InventoryItemID {
	// 				return true
	// 			}
	// 		}
	// 		return false
	// 	}(vendorItem)
	// 	if !contains {
	// 		vendorItems = append(vendorItems, vendorItem)
	// 	}
	// }

	// return &vendorItems
	return nil
}

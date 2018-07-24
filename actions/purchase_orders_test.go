package actions

import (
	"buddhabowls/logic"
	"buddhabowls/models"
	"buddhabowls/presentation"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"net/url"
	"time"
)

var _ = fmt.Printf

// viewing the purchase orders at a certain week
func (as *ActionSuite) Test_ListPO_View() {
	orderTime := time.Date(2018, 7, 5, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)

	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", orderTime.Format(time.RFC3339))).Get()
	as.Contains(res.Body.String(), purchaseOrder.ID.String())
}

// viewing no PO's
func (as *ActionSuite) Test_ListPO_NoResults() {
	orderTime := time.Date(2018, 7, 5, 0, 0, 0, 0, time.UTC)

	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", orderTime.Format(time.RFC3339))).Get()
	as.NotContains(res.Body.String(), "Open Orders")
	as.NotContains(res.Body.String(), "Received Orders")
}

// not viewing purchase order when out of week
func (as *ActionSuite) Test_ListPO_OutOfWeek() {
	orderTime := time.Date(2018, 7, 5, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)

	orderTime = orderTime.AddDate(0, 0, -7)
	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", orderTime.Format(time.RFC3339))).Get()
	as.NotContains(res.Body.String(), purchaseOrder.ID.String())
	as.NotContains(res.Body.String(), "Open Orders")
	as.NotContains(res.Body.String(), "Received Orders")
}

// viewing PO's in a custom date range
func (as *ActionSuite) Test_ListPO_CustomDate() {
	orderTime := time.Date(2018, 7, 5, 0, 0, 0, 0, time.UTC)
	earlyOrderTime := time.Date(2018, 6, 15, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)
	otherPO, err := createPO(as.DB, earlyOrderTime)
	as.NoError(err)

	orderTime = orderTime.AddDate(0, 0, 1)
	earlyOrderTime = earlyOrderTime.AddDate(0, 0, -1)
	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s&EndTime=%s", earlyOrderTime.Format(time.RFC3339), orderTime.Format(time.RFC3339))).Get()

	as.Contains(res.Body.String(), purchaseOrder.ID.String())
	as.Contains(res.Body.String(), otherPO.ID.String())
	as.Contains(res.Body.String(), "Open Orders")
	as.NotContains(res.Body.String(), "Received Orders")
}

// viewing PO created at the last second of the week
func (as *ActionSuite) Test_ListPO_LastSecond() {
	// week ends 7/8/2018
	orderTime := time.Date(2018, 7, 9, 0, 0, 0, 0, time.UTC).Add(-time.Second)
	startTime := time.Date(2018, 7, 2, 0, 0, 0, 0, time.UTC)

	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)

	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", startTime.Format(time.RFC3339))).Get()

	as.Contains(res.Body.String(), purchaseOrder.ID.String())
	as.Contains(res.Body.String(), "Open Orders")
	as.NotContains(res.Body.String(), "Received Orders")
}

// viewing PO created at the first second of the week
func (as *ActionSuite) Test_ListPO_FirstSecond() {
	// week starts 7/2/2018
	startTime := time.Date(2018, 7, 2, 0, 0, 0, 0, time.UTC)

	purchaseOrder, err := createPO(as.DB, startTime)
	as.NoError(err)

	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", startTime.Format(time.RFC3339))).Get()

	as.Contains(res.Body.String(), purchaseOrder.ID.String())
	as.Contains(res.Body.String(), "Open Orders")
	as.NotContains(res.Body.String(), "Received Orders")
}

// viewing only received orders
func (as *ActionSuite) Test_ListPO_Received() {
	orderTime := time.Date(2018, 7, 5, 0, 0, 0, 0, time.UTC)

	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)
	as.NoError(receiveOrder(as.DB, purchaseOrder, orderTime))

	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", orderTime.Format(time.RFC3339))).Get()

	as.Contains(res.Body.String(), purchaseOrder.ID.String())
	as.NotContains(res.Body.String(), "Open Orders")
	as.Contains(res.Body.String(), "Received Orders")
}

// viewing both open and received
func (as *ActionSuite) Test_ListPO_OpenAndReceived() {
	orderTime := time.Date(2018, 7, 5, 0, 0, 0, 0, time.UTC)

	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)
	recPurchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)
	as.NoError(receiveOrder(as.DB, recPurchaseOrder, orderTime))

	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", orderTime.Format(time.RFC3339))).Get()

	as.Contains(res.Body.String(), purchaseOrder.ID.String())
	as.Contains(res.Body.String(), recPurchaseOrder.ID.String())
	as.Contains(res.Body.String(), "Open Orders")
	as.Contains(res.Body.String(), "Received Orders")
}

// viewing new PO page
func (as *ActionSuite) Test_NewPO_View() {
	vendor, err := createVendor(as.DB)
	as.NoError(err)

	res := as.HTML("/purchase_orders/new").Get()

	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), vendor.ID.String())
	for _, item := range vendor.Items {
		as.Contains(res.Body.String(), item.ID.String())
		as.Contains(res.Body.String(), item.InventoryItemID.String())
	}
}

// creating PO with all items count of 1
func (as *ActionSuite) Test_CreatePO() {
	vendor, err := createVendor(as.DB)
	as.NoError(err)

	items := presentation.NewItemsAPI(vendor.Items)
	for i := 0; i < len(items); i++ {
		items[i].Count = 1
	}

	orderDate := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	itemsJSON, err := json.Marshal(items)
	as.NoError(err)
	formData := struct {
		OrderDate time.Time
		VendorID  string
		Items     string
	}{
		orderDate,
		vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML("/purchase_orders").Post(formData)

	as.Equal(303, res.Code)

	resultURL, err := url.Parse(res.Location())
	as.NoError(err)
	path := resultURL.EscapedPath()
	as.Equal("/purchase_orders", path)

	purchaseOrders, err := logic.GetPurchaseOrders(orderDate, orderDate, as.DB)
	as.NoError(err)
	purchaseOrder := (*purchaseOrders)[0]

	as.Equal(orderDate.Unix(), purchaseOrder.OrderDate.Time.Unix())
	as.Equal(vendor.Items[0].InventoryItemID.String(), purchaseOrder.Items[0].InventoryItemID.String())
	as.NotEqual(vendor.Items[0].ID.String(), purchaseOrder.Items[0].ID.String())
	as.NotEqual(uuid.UUID{}.String(), purchaseOrder.Items[0].ID.String())
	as.Equal(1.0, purchaseOrder.Items[0].Count)
}

// creating PO with no items (produce error message)
func (as *ActionSuite) Test_CreatePO_NoItemCounts() {
	vendor, err := createVendor(as.DB)
	as.NoError(err)

	items := presentation.NewItemsAPI(vendor.Items)
	for i := 0; i < len(items); i++ {
		items[i].Count = 0
	}

	orderDate := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	itemsJSON, err := json.Marshal(items)
	as.NoError(err)
	formData := struct {
		OrderDate time.Time
		VendorID  string
		Items     string
	}{
		orderDate,
		vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML("/purchase_orders").Post(formData)

	as.Equal(422, res.Code)

	purchaseOrder := &models.PurchaseOrder{}
	as.Error(as.DB.First(purchaseOrder))
}

// viewing edit page for PO
func (as *ActionSuite) Test_EditPO_View() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s/edit", purchaseOrder.ID.String())).Get()

	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "07/04/2018")
	for _, item := range purchaseOrder.Items {
		as.Contains(res.Body.String(), item.ID.String())
		as.Contains(res.Body.String(), item.InventoryItemID.String())
	}
}

// viewing edit page for received PO (fill in later)

// editing PO changing counts
func (as *ActionSuite) Test_EditPO_ChangeCounts() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)

	items := presentation.NewItemsAPI(purchaseOrder.Items)
	for i := 0; i < len(items); i++ {
		items[i].Count = 69
	}

	newOrderTime := time.Date(2018, 7, 6, 0, 0, 0, 0, time.UTC)
	itemsJSON, err := json.Marshal(items)
	as.NoError(err)
	formData := struct {
		OrderDate time.Time
		VendorID  string
		Items     string
	}{
		newOrderTime,
		purchaseOrder.Vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s", purchaseOrder.ID.String())).Put(formData)

	// ensure redirect
	as.Equal(303, res.Code)

	// ensure redirect to the index page
	resultURL, err := url.Parse(res.Location())
	as.NoError(err)
	path := resultURL.EscapedPath()
	as.Equal("/purchase_orders", path)

	dbPOs, err := logic.GetPurchaseOrders(newOrderTime, newOrderTime, as.DB)
	as.NoError(err)

	dbPO := (*dbPOs)[0]
	for _, item := range dbPO.Items {
		as.Equal(69.0, item.Count)
	}
}

// editing PO adding items
func (as *ActionSuite) Test_EditPO_AddItem() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)

	items := presentation.NewItemsAPI(purchaseOrder.Items)
	newVendorItem, err := createVendorItem(as.DB, &purchaseOrder.Vendor, "another_item")
	as.NoError(err)

	newItem := presentation.NewItemAPI(*newVendorItem)
	items = append(items, newItem)
	for i := 0; i < len(items); i++ {
		items[i].Count = 69
	}

	newOrderTime := time.Date(2018, 7, 6, 0, 0, 0, 0, time.UTC)
	itemsJSON, err := json.Marshal(items)
	as.NoError(err)
	formData := struct {
		OrderDate time.Time
		VendorID  string
		Items     string
	}{
		newOrderTime,
		purchaseOrder.Vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s", purchaseOrder.ID.String())).Put(formData)

	// ensure redirect
	as.Equal(303, res.Code)

	dbPOs, err := logic.GetPurchaseOrders(newOrderTime, newOrderTime, as.DB)
	as.NoError(err)

	dbPO := (*dbPOs)[0]
	for _, item := range dbPO.Items {
		as.Equal(69.0, item.Count)
	}
	as.Equal(2, len(dbPO.Items))
}

// editing PO removing items
func (as *ActionSuite) Test_EditPO_RemoveItem() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPOMultipleItems(as.DB, orderTime)
	as.NoError(err)

	for _, item := range purchaseOrder.Items {
		if item.GetName() == "yet_another_item" {
			purchaseOrder.Items = models.OrderItems{item}
			break
		}
	}

	items := presentation.NewItemsAPI(purchaseOrder.Items)
	itemsJSON, err := json.Marshal(items)
	as.NoError(err)
	formData := struct {
		OrderDate time.Time
		VendorID  string
		Items     string
	}{
		orderTime,
		purchaseOrder.Vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s", purchaseOrder.ID.String())).Put(formData)

	// ensure redirect
	as.Equal(303, res.Code)

	dbPOs, err := logic.GetPurchaseOrders(orderTime, orderTime, as.DB)
	as.NoError(err)

	dbPO := (*dbPOs)[0]
	as.Equal(1, len(dbPO.Items))
	as.Equal("yet_another_item", dbPO.Items[0].GetName())
}

// editing PO removing all items (produces error)
func (as *ActionSuite) Test_EditPO_RemoveAllItem() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPOMultipleItems(as.DB, orderTime)
	as.NoError(err)

	for i := 0; i < len(purchaseOrder.Items); i++ {
		purchaseOrder.Items[i].Count = 0
	}

	items := presentation.NewItemsAPI(purchaseOrder.Items)
	itemsJSON, err := json.Marshal(items)
	as.NoError(err)
	formData := struct {
		OrderDate time.Time
		VendorID  string
		Items     string
	}{
		orderTime,
		purchaseOrder.Vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s", purchaseOrder.ID.String())).Put(formData)

	// ensure error code
	as.Equal(422, res.Code)

	presenter := presentation.NewPresenter(as.DB)
	dbPOs, err := presenter.GetPurchaseOrders(orderTime, orderTime)
	as.NoError(err)
	dbPO := (*dbPOs)[0]
	as.Equal(2, len(dbPO.Items))
}

// editing open PO, setting to received
func (as *ActionSuite) Test_EditPO_SetReceived() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)

	receivedTime := time.Date(2018, 7, 6, 0, 0, 0, 0, time.UTC)
	itemsJSON, err := json.Marshal(presentation.NewItemsAPI(purchaseOrder.Items))
	as.NoError(err)
	formData := struct {
		OrderDate    time.Time
		ReceivedDate time.Time
		VendorID     string
		Items        string
	}{
		orderTime,
		receivedTime,
		purchaseOrder.Vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s", purchaseOrder.ID.String())).Put(formData)

	// ensure redirect
	as.Equal(303, res.Code)

	// ensure redirect to the index page
	resultURL, err := url.Parse(res.Location())
	as.NoError(err)
	path := resultURL.EscapedPath()
	as.Equal("/purchase_orders", path)

	dbPOs, err := logic.GetPurchaseOrders(orderTime, orderTime, as.DB)
	as.NoError(err)

	dbPO := (*dbPOs)[0]
	as.Equal(receivedTime.Unix(), dbPO.ReceivedDate.Time.Unix())
}

// editing received PO, setting to open
func (as *ActionSuite) Test_EditPO_SetReOpen() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)
	receivedTime := time.Date(2018, 7, 6, 0, 0, 0, 0, time.UTC)
	purchaseOrder.ReceivedDate.Time = receivedTime
	purchaseOrder.ReceivedDate.Valid = true
	as.DB.Update(purchaseOrder)

	itemsJSON, err := json.Marshal(presentation.NewItemsAPI(purchaseOrder.Items))
	as.NoError(err)
	formData := struct {
		OrderDate time.Time
		VendorID  string
		Items     string
	}{
		orderTime,
		purchaseOrder.Vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s", purchaseOrder.ID.String())).Put(formData)

	// ensure redirect
	as.Equal(303, res.Code)

	// ensure redirect to the index page
	resultURL, err := url.Parse(res.Location())
	as.NoError(err)
	path := resultURL.EscapedPath()
	as.Equal("/purchase_orders", path)

	dbPOs, err := logic.GetPurchaseOrders(orderTime, orderTime, as.DB)
	as.NoError(err)

	dbPO := (*dbPOs)[0]
	as.False(dbPO.ReceivedDate.Valid)
}

// editing PO, setting received before open date (produces error)
func (as *ActionSuite) Test_EditPO_SetReceivedError() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)
	receivedTime := time.Date(2018, 7, 3, 0, 0, 0, 0, time.UTC)
	purchaseOrder.ReceivedDate.Time = receivedTime
	purchaseOrder.ReceivedDate.Valid = true

	itemsJSON, err := json.Marshal(presentation.NewItemsAPI(purchaseOrder.Items))
	as.NoError(err)
	formData := struct {
		OrderDate    time.Time
		ReceivedDate time.Time
		VendorID     string
		Items        string
	}{
		orderTime,
		receivedTime,
		purchaseOrder.Vendor.ID.String(),
		string(itemsJSON),
	}

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s", purchaseOrder.ID.String())).Put(formData)

	as.Equal(422, res.Code)

	as.Contains(res.Body.String(), "error")

	dbPOs, err := logic.GetPurchaseOrders(orderTime, orderTime, as.DB)
	as.NoError(err)

	dbPO := (*dbPOs)[0]
	as.False(dbPO.ReceivedDate.Valid)
}

// receiving open PO

// receiving open PO whose order date is in the future (produces error)

// re-opening received PO

// changing order date on open PO

// changing order date on received PO
// (should user be able to do this?)

// changing received date on received PO

// changing received date to earlier than order date (produces error)

// changing order date to after received date (produces error)

// destroying PO
func (as *ActionSuite) Test_DestroyPO() {
	orderTime := time.Date(2018, 7, 4, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)

	res := as.HTML(fmt.Sprintf("/purchase_orders/%s", purchaseOrder.ID.String())).Delete()

	as.Equal(303, res.Code)

	// ensure redirect to the index page
	resultURL, err := url.Parse(res.Location())
	as.NoError(err)
	path := resultURL.EscapedPath()
	as.Equal("/purchase_orders", path)

	dbPOs, err := logic.GetPurchaseOrders(orderTime, orderTime, as.DB)
	as.NoError(err)

	as.Equal(0, len(*dbPOs))
}

func createPO(db *pop.Connection, orderTime time.Time) (*models.PurchaseOrder, error) {
	vendor, err := createVendor(db)
	if err != nil {
		return nil, err
	}

	purchaseOrder := &models.PurchaseOrder{
		Vendor:    *vendor,
		VendorID:  vendor.ID,
		OrderDate: nulls.Time{Time: orderTime, Valid: true},
	}
	if err = db.Create(purchaseOrder); err != nil {
		return nil, err
	}

	item := &models.OrderItem{
		InventoryItem:   vendor.Items[0].InventoryItem,
		InventoryItemID: vendor.Items[0].InventoryItemID,
		Count:           4,
		Price:           vendor.Items[0].Price,
		OrderID:         purchaseOrder.ID,
	}

	if err = db.Create(item); err != nil {
		return nil, err
	}

	purchaseOrder.Items = models.OrderItems{*item}

	return purchaseOrder, nil
}

func createPOMultipleItems(db *pop.Connection, orderTime time.Time) (*models.PurchaseOrder, error) {
	purchaseOrder, err := createPO(db, orderTime)
	if err != nil {
		return nil, err
	}

	newItem, err := createVendorItem(db, &purchaseOrder.Vendor, "yet_another_item")
	item := &models.OrderItem{
		InventoryItem:   newItem.InventoryItem,
		InventoryItemID: newItem.InventoryItemID,
		Count:           6,
		Price:           newItem.Price,
		OrderID:         purchaseOrder.ID,
	}

	if err = db.Create(item); err != nil {
		return nil, err
	}

	purchaseOrder.Items = append(purchaseOrder.Items, *item)

	return purchaseOrder, nil
}

func receiveOrder(db *pop.Connection, purchaseOrder *models.PurchaseOrder, date time.Time) error {
	purchaseOrder.ReceivedDate = nulls.Time{Time: date, Valid: true}
	return db.Update(purchaseOrder)
}

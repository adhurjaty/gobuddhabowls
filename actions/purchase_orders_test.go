package actions

import (
	"buddhabowls/models"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"time"
)

var _ = fmt.Printf

// viewing the purchase orders at a certain week
func (as *ActionSuite) Test_ListPO() {
	orderTime := time.Date(2018, 7, 5, 0, 0, 0, 0, time.UTC)
	purchaseOrder, err := createPO(as.DB, orderTime)
	as.NoError(err)

	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", orderTime.Format(time.RFC3339))).Get()
	as.Contains(res.Body.String(), purchaseOrder.ID.String())
}

// viewing no PO's
func (as *ActionSuite) Test_ListPONoResults() {
	orderTime := time.Date(2018, 7, 5, 0, 0, 0, 0, time.UTC)

	res := as.HTML(fmt.Sprintf("/purchase_orders?StartTime=%s", orderTime.Format(time.RFC3339))).Get()
	as.NotContains(res.Body.String(), "Open Orders")
	as.NotContains(res.Body.String(), "Received Orders")
}

// not viewing purchase order when out of week
func (as *ActionSuite) Test_ListPOOutOfWeek() {
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
func (as *ActionSuite) Test_ListPOCustomDate() {
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
func (as *ActionSuite) Test_ListPOLastSecond() {
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
// viewing PO created at the last second of the week
func (as *ActionSuite) Test_ListPOFirstSecond() {
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

// viewing both open and received

// viewing new PO page

// creating PO with all items count of 1

// creating PO with some items

// creating PO with no items (produce error message)

// viewing edit page for PO

// editing PO changing counts

// editing PO adding items

// editing PO removing items

// editing PO removing all items (produces error)

// editing open PO, setting to received

// editing received PO, setting to open

// editing PO, setting received before open date (produces error)

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

// func (as *ActionSuite) Test_PurchaseOrdersResource_Show() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_PurchaseOrdersResource_New() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_PurchaseOrdersResource_Create() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_PurchaseOrdersResource_Edit() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_PurchaseOrdersResource_Update() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_PurchaseOrdersResource_Destroy() {
// 	as.Fail("Not Implemented!")
// }

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

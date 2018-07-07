package actions

import (
	"buddhabowls/models"
	"database/sql"
	"github.com/gobuffalo/pop"
)

func (as *ActionSuite) Test_VendorItemsResource_List() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorItemsResource_Show() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorItemsResource_New() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorItemsResource_Create() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorItemsResource_Edit() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorItemsResource_Update() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorItemsResource_Destroy() {
	as.Fail("Not Implemented!")
}

func createVendorItem(db *pop.Connection, vendor *models.Vendor, name string) (*models.VendorItem, error) {
	invItem, err := createInventoryItem(db, name)
	if err != nil {
		return nil, err
	}

	item := &models.VendorItem{
		InventoryItem:   *invItem,
		InventoryItemID: invItem.ID,
		Vendor:          *vendor,
		VendorID:        vendor.ID,
		Conversion:      4,
		Price:           8.99,
		PurchasedUnit:   sql.NullString{String: "Quart", Valid: true},
	}
	if err = db.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

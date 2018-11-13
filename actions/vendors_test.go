package actions

import (
	"buddhabowls/models"
	"database/sql"
	"github.com/gobuffalo/pop"
)

func (as *ActionSuite) Test_ListVendor_View() {
	vendor, err := createVendor(as.DB)
	as.NoError(err)

	res := as.HTML("/vendors").Get()
	as.Contains(res.Body.String(), vendor.ID.String())
}

func (as *ActionSuite) Test_VendorsResource_Show() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_New() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_Create() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_Edit() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_Update() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_Destroy() {
	as.Fail("Not Implemented!")
}

func createVendor(db *pop.Connection) (*models.Vendor, error) {
	vendor := &models.Vendor{
		Name:         "test_vendor",
		ShippingCost: 2.5,
	}

	err := db.Create(vendor)
	if err != nil {
		return nil, err
	}

	invItem, err := createInventoryItem(db, "test_item")
	if err != nil {
		return nil, err
	}

	item := &models.VendorItem{
		InventoryItem:   *invItem,
		InventoryItemID: invItem.ID,
		Vendor:          *vendor,
		VendorID:        vendor.ID,
		Conversion:      1,
		Price:           5,
		PurchasedUnit:   sql.NullString{String: "EA", Valid: true},
	}
	if err = db.Create(item); err != nil {
		return nil, err
	}

	vendor.Items = models.VendorItems{*item}

	return vendor, nil
}

package actions

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
)

func (as *ActionSuite) Test_InventoryItemsResource_List() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_InventoryItemsResource_Show() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_InventoryItemsResource_New() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_InventoryItemsResource_Create() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_InventoryItemsResource_Edit() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_InventoryItemsResource_Update() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_InventoryItemsResource_Destroy() {
	as.Fail("Not Implemented!")
}

func createInventoryItem(db *pop.Connection) (*models.InventoryItem, error) {
	category := &models.InventoryItemCategory{
		Name:  "test_category",
		Index: 0,
	}
	err := db.Create(category)
	if err != nil {
		return nil, err
	}
	item := &models.InventoryItem{
		Name:       "test_item",
		Category:   *category,
		CategoryID: category.ID,
		IsActive:   true,
		Index:      0,
	}
	err = db.Create(item)

	return item, err
}

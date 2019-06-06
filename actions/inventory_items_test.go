package actions

import (
	"buddhabowls/models"
	"buddhabowls/presentation"
	"fmt"
)

var _ = fmt.Println

type FormVendorItem struct {
	Conversion    float64
	PurchasedUnit string
	Price         float64
}

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
	itemName := "TestItem"

	category := &models.ItemCategory{
		Name:  "Dairy",
		Index: 1,
	}

	as.DB.Create(category)

	vendors := &models.Vendors{
		models.Vendor{
			Name: "ThisVendor",
		},
		models.Vendor{
			Name: "ThatVendor",
		},
	}

	vMap := make(map[string]FormVendorItem)
	for i, vendor := range *vendors {
		as.DB.Create(&(*vendors)[i])
		vMap[vendor.Name] = FormVendorItem{
			Conversion:    1,
			PurchasedUnit: "EA",
			Price:         4.2,
		}
	}

	login(as.DB, as.Session)

	catAPI := presentation.NewCategoryAPI(*category)
	formData := struct {
		Category             *presentation.CategoryAPI
		VendorItemMap        map[string]FormVendorItem
		Name                 string
		CountUnit            string
		RecipeUnit           string
		RecipeUnitConversion float64
		Yield                float64
		Index                int
	}{
		&catAPI,
		vMap,
		itemName,
		"EA",
		"RU",
		1,
		1,
		1,
	}
	res := as.HTML("/inventory_items").Post(formData)

	as.Equal(303, res.Code)

	dbItem := &models.InventoryItem{}
	err := as.DB.Where("name = ?", itemName).First(dbItem)
	as.NoError(err)
	as.Equal(itemName, dbItem.Name)

	vendItems := &models.VendorItems{}
	err = as.DB.All(vendItems)
	as.NoError(err)
	as.Equal(2, len(*vendItems))
	for _, item := range *vendItems {
		as.True(item.VendorID.String() == (*vendors)[0].ID.String() ||
			item.VendorID.String() == (*vendors)[1].ID.String())
	}
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

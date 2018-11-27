package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"github.com/gobuffalo/uuid"
	"time"
)

type InventoryAPI struct {
	ID    uuid.UUID `json:"id"`
	Date  time.Time `json:"time"`
	Items ItemsAPI  `json:"Items"`
}

type InventoriesAPI []InventoryAPI

func (inv InventoryAPI) String() string {
	jo, _ := json.Marshal(inv)
	return string(jo)
}

func (inv InventoriesAPI) String() string {
	jo, _ := json.Marshal(inv)
	return string(jo)
}

func NewInventoryAPI(inventory *models.Inventory, vendors *VendorsAPI) InventoryAPI {
	items := NewItemsAPI(inventory.Items)
	populateVendorItems(&items, &inventory.Items, vendors)

	return InventoryAPI{
		ID:    inventory.ID,
		Date:  inventory.Date,
		Items: NewItemsAPI(inventory.Items),
	}
}

func NewInventoriesAPI(inventories *models.Inventories, vendors *VendorsAPI) InventoriesAPI {
	apis := make([]InventoryAPI, len(*inventories))
	for i, inventory := range *inventories {
		apis[i] = NewInventoryAPI(&inventory, vendors)
	}

	return apis
}

func populateVendorItems(apiItems *ItemsAPI, items *models.CountInventoryItems, vendors *VendorsAPI) {
	for i, item := range *items {
		vendorMap := make(map[string]ItemAPI)
		for _, vendor := range *vendors {
			vendorItem := getVendorItem(item.InventoryItemID, &vendor.Items)
			if vendorItem != nil {
				vendorMap[vendor.Name] = *vendorItem
			}
		}
		(*apiItems)[i].VendorItemMap = vendorMap
		(*apiItems)[i].SetSelectedVendor(item.SelectedVendor.Name)
	}
}

func getVendorItem(inventoryItemID uuid.UUID, items *ItemsAPI) *ItemAPI {
	for _, item := range *items {
		if inventoryItemID.String() == item.InventoryItemID {
			return &item
		}
	}

	return nil
}

// func getVendorName(id uuid.NullUUID, vendors *VendorsAPI) string {
// 	idVal, err := id.Value()
// 	if err != nil {
// 		return ""
// 	}
// 	for _, vendor := range *vendors {
// 		if idVal.(uuid.UUID) == vendor.ID {
// 			return vendor.Name
// 		}
// 	}

// 	return ""
// }

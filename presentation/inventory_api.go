package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"github.com/gobuffalo/uuid"
	"time"
)

type InventoryAPI struct {
	ID    string    `json:"id"`
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
		ID:    inventory.ID.String(),
		Date:  inventory.Date,
		Items: items,
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
	if vendors == nil {
		return
	}

	for i, item := range *items {
		(*apiItems)[i].VendorItemMap = GetVendorMap(item.InventoryItemID.String(),
			vendors)
		(*apiItems)[i].SetSelectedVendor(item.SelectedVendor.Name)
	}
}

func GetVendorMap(invItemID string, vendors *VendorsAPI) map[string]ItemAPI {
	vendorMap := make(map[string]ItemAPI)
	for _, vendor := range *vendors {
		vendorItem := getVendorItem(invItemID, &vendor.Items)
		if vendorItem != nil {
			// HACK: getting vendor ID from the property SelectedVendor within
			// the vendor item map
			vendorItem.SelectedVendor = vendor.ID
			vendorMap[vendor.Name] = *vendorItem
		}
	}

	return vendorMap
}

func getVendorItem(inventoryItemID string, items *ItemsAPI) *ItemAPI {
	for _, item := range *items {
		if inventoryItemID == item.InventoryItemID {
			return &item
		}
	}

	return nil
}

func ConvertToModelInventory(invAPI *InventoryAPI) (*models.Inventory, error) {
	id := uuid.UUID{}
	if len(invAPI.ID) > 0 {
		var err error
		id, err = uuid.FromString(invAPI.ID)
		if err != nil {
			return nil, err
		}
	}

	items, err := ConvertToModelCountInventoryItems(invAPI.Items, id)
	if err != nil {
		return nil, err
	}

	return &models.Inventory{
		ID:    id,
		Date:  invAPI.Date,
		Items: *items,
	}, nil
}

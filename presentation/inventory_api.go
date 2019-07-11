package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/uuid"
)

type InventoryAPI struct {
	ID        string    `json:"id"`
	Date      time.Time `json:"time"`
	Items     ItemsAPI  `json:"Items"`
	PrepItems ItemsAPI  `json:"PrepItems"`
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

func NewInventoryAPI(inventory *models.Inventory, vendors *VendorsAPI, presenter *Presenter) (InventoryAPI, error) {
	items := NewItemsAPI(&inventory.Items)
	populateVendorItems(&items, &inventory.Items, vendors)

	prepItems := NewItemsAPI(&inventory.PrepItems)
	err := presenter.populatePrepItemCosts(&prepItems)
	fmt.Println(prepItems)

	return InventoryAPI{
		ID:        inventory.ID.String(),
		Date:      inventory.Date,
		Items:     items,
		PrepItems: prepItems,
	}, err
}

func NewInventoriesAPI(inventories *models.Inventories, presenter *Presenter) (InventoriesAPI, error) {
	vendors, err := presenter.GetVendors()
	if err != nil {
		return InventoriesAPI{}, err
	}

	apis := make([]InventoryAPI, len(*inventories))
	for i, inventory := range *inventories {
		invAPI, err := NewInventoryAPI(&inventory, vendors, presenter)
		if err != nil {
			return apis, err
		}

		apis[i] = invAPI
	}

	return apis, nil
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

	prepItems, err := ConvertToModelCountPrepItems(invAPI.PrepItems, id)
	if err != nil {
		return nil, err
	}

	return &models.Inventory{
		ID:        id,
		Date:      invAPI.Date,
		Items:     *items,
		PrepItems: *prepItems,
	}, nil
}

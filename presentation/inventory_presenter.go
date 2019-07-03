package presentation

import (
	"buddhabowls/helpers"
	"buddhabowls/logic"
	"buddhabowls/models"
	"time"

	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

func (p *Presenter) GetInventories(startTime time.Time, endTime time.Time) (*InventoriesAPI, error) {
	inventories, err := logic.GetInventories(startTime, endTime, p.tx)
	if err != nil {
		return nil, err
	}

	vendors, err := p.GetVendors()
	if err != nil {
		return nil, err
	}

	apiInv := NewInventoriesAPI(inventories, vendors)

	return &apiInv, nil
}

func (p *Presenter) GetLatestInventory(date time.Time) (*InventoryAPI, error) {
	inventory, err := logic.GetLatestInventory(date, p.tx)
	if err != nil {
		return nil, err
	}

	vendors, err := p.GetVendors()
	if err != nil {
		return nil, err
	}

	apiInv := NewInventoryAPI(inventory, vendors)
	return &apiInv, nil
}

func (p *Presenter) GetInventory(id string) (*InventoryAPI, error) {
	inventory, err := logic.GetInventory(id, p.tx)
	if err != nil {
		return nil, err
	}

	apiInv := NewInventoryAPI(inventory, nil)
	return &apiInv, nil
}

func (p *Presenter) InsertInventory(invAPI *InventoryAPI) (*validate.Errors, error) {
	inventory, err := ConvertToModelInventory(invAPI)
	if err != nil {
		return nil, err
	}

	verrs, err := p.updateInvVendorItems(invAPI.Items)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	return logic.InsertInventory(inventory, p.tx)
}

func (p *Presenter) UpdateInventory(invAPI *InventoryAPI) (*validate.Errors, error) {
	inventory, err := ConvertToModelInventory(invAPI)
	if err != nil {
		return nil, err
	}

	if invAPI.Items != nil && len(invAPI.Items) > 0 {
		verrs, err := p.updateInvVendorItems(invAPI.Items)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return logic.UpdateInventory(inventory, p.tx)
}

func (p *Presenter) updateInvVendorItems(invItems ItemsAPI) (*validate.Errors, error) {
	vendorItems := models.VendorItems{}
	for _, item := range invItems {
		subItem, ok := item.VendorItemMap[item.SelectedVendor]
		if ok {
			vendorID, err := uuid.FromString(subItem.SelectedVendor)
			if err != nil {
				continue
			}
			vItem, err := ConvertToModelVendorItem(subItem, vendorID)
			if err != nil {
				continue
			}

			vendorItems = append(vendorItems, *vItem)
		}
	}

	return logic.UpdateVendorItems(&vendorItems, p.tx)
}

func (p *Presenter) DestroyInventory(invAPI *InventoryAPI) error {
	inventory, err := ConvertToModelInventory(invAPI)
	if err != nil {
		return err
	}

	return logic.DestroyInventory(inventory, p.tx)
}

// func (p *Presenter) getCountInventoryItem(item *ItemAPI) (*models.CountInventoryItem, error) {
// 	inventoryID, err := p.getLatestInventoryID(item)
// 	if err != nil {
// 		return nil, err
// 	}
// 	countInvItem, err := ConvertToModelCountInventoryItem(*item, inventoryID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dbItem, err := logic.GetCountInventoryItemByInvItem(item.InventoryItemID,
// 		inventoryID.String(), p.tx)
// 	if err != nil {
// 		countInvItem.InventoryID = inventoryID
// 		countInvItem.Count = 0
// 		_, _ = logic.InsertCountInventoryItem(countInvItem, p.tx)

// 		dbItem, err = logic.GetCountInventoryItemByInvItem(item.InventoryItemID,
// 			inventoryID.String(), p.tx)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	countInvItem.ID = dbItem.ID
// 	countInvItem.InventoryID = dbItem.InventoryID

// 	return countInvItem, nil
// }

func (p *Presenter) getLatestInventoryID(item *ItemAPI) (uuid.UUID, error) {
	inventory, err := logic.GetLatestInventory(helpers.Today(), p.tx)
	if err != nil {
		return uuid.UUID{}, err
	}

	return inventory.ID, nil
}

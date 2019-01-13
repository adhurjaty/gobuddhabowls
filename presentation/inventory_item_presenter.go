package presentation

import (
	"buddhabowls/logic"
	"fmt"
	"github.com/gobuffalo/validate"
	"time"
)

var _ = fmt.Println

func (p *Presenter) GetInventoryItems() (*ItemsAPI, error) {
	items, err := logic.GetInventoryItems(p.tx)
	if err != nil {
		return nil, err
	}

	apiItems := NewItemsAPI(*items)

	return &apiItems, err
}

func (p *Presenter) GetMasterInventoryList() (*ItemsAPI, error) {
	items, err := p.GetNewInventoryItems()
	if err != nil {
		return nil, err
	}

	if err = p.populateInventoryItemDetails(items); err != nil {
		return nil, err
	}

	return items, nil
}

func (p *Presenter) populateInventoryItemDetails(items *ItemsAPI) error {
	for i := 0; i < len(*items); i++ {
		item := &(*items)[i]
		invItem, err := logic.GetInventoryItem(item.InventoryItemID, p.tx)
		if err != nil {
			return err
		}

		item.ID = invItem.ID.String()
		item.RecipeUnit = invItem.RecipeUnit
		item.RecipeUnitConversion = invItem.RecipeUnitConversion
	}

	return nil
}

func (p *Presenter) GetNewInventoryItems() (*ItemsAPI, error) {
	// get base inventory items
	items, err := p.GetInventoryItems()
	if err != nil {
		return nil, err
	}

	// populate them based on latest inventory
	if err = p.populateLatestInvItems(items); err != nil {
		return nil, err
	}
	// populate the latest selected vendor
	p.populateSelectedVendors(items)

	clearItemIds(items)

	return items, nil
}

func (p *Presenter) populateLatestInvItems(items *ItemsAPI) error {
	latestInv, err := p.GetLatestInventory(time.Now())
	if err != nil {
		return err
	}

	var vendors *VendorsAPI

	for i := 0; i < len(*items); i++ {
		item := &(*items)[i]
		for _, latestItem := range latestInv.Items {
			if item.InventoryItemID == latestItem.InventoryItemID {
				item.Count = latestItem.Count
				item.VendorItemMap = latestItem.VendorItemMap

				// default behavior, will probably be re-set in the next function
				item.SetSelectedVendor(latestItem.SelectedVendor)
				break
			}
		}
		if item.VendorItemMap == nil {
			item.Count = 0
			if vendors == nil {
				vendors, err = p.GetVendors()
				if err != nil {
					return err
				}
			}
			item.VendorItemMap = GetVendorMap(item.InventoryItemID, vendors)
		}
	}

	return nil
}

func (p *Presenter) populateSelectedVendors(items *ItemsAPI) {
	for i := 0; i < len(*items); i++ {
		item := &(*items)[i]
		vendor, err := logic.GetLatestVendor(item.InventoryItemID, p.tx)
		if err != nil {
			continue
		}

		item.SetSelectedVendor(vendor.Name)
	}
}

func clearItemIds(items *ItemsAPI) {
	for _, item := range *items {
		item.ID = ""
	}
}

func (p *Presenter) GetInventoryItem(id string) (*ItemAPI, error) {
	item, err := logic.GetInventoryItem(id, p.tx)
	if err != nil {
		return nil, err
	}

	apiItem := NewItemAPI(item)
	return &apiItem, nil
}

func (p *Presenter) UpdateInventoryItem(item *ItemAPI) (*validate.Errors, error) {
	invItem, err := ConvertToModelInventoryItem(item)
	if err != nil {
		return nil, err
	}

	vendorItem, err := p.getVendorItem(item)
	if err != nil {
		return nil, err
	}

	selectedItem, _ := item.VendorItemMap[item.SelectedVendor]
	(&selectedItem).SelectedVendor = vendorItem.VendorID.String()

	verrs, err := logic.UpdateInventoryItem(invItem, p.tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return logic.UpdateVendorItem(vendorItem, p.tx)
}

func (p *Presenter) InsertInventoryItem(item *ItemAPI) (*validate.Errors, error) {
	invItem, err := ConvertToModelInventoryItem(item)
	if err != nil {
		return nil, err
	}

	verrs, err := logic.InsertInventoryItem(invItem, p.tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	for selectedVendor, vItemAPI := range item.VendorItemMap {
		vendorID, err := p.getVendorID(selectedVendor)
		if err != nil {
			return nil, err
		}
		vItemAPI.InventoryItemID = invItem.ID.String()
		vendorItem, err := ConvertToModelVendorItem(vItemAPI, vendorID)
		if err != nil {
			return nil, err
		}

		vendorItem.VendorID = vendorID
		vendorItem.InventoryItemID = invItem.ID

		verrs, err = logic.InsertVendorItem(vendorItem, p.tx)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return verrs, nil
}

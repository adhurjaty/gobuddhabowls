package presentation

import (
	"buddhabowls/logic"
	"fmt"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"strings"
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

	for i := 0; i < len(*items); i++ {
		item := &(*items)[i]
		if err = p.populateInventoryItemDetails(item); err != nil {
			return nil, err
		}
	}

	return items, nil
}

func (p *Presenter) populateInventoryItemDetails(item *ItemAPI) error {
	invItem, err := logic.GetInventoryItem(item.InventoryItemID, p.tx)
	if err != nil {
		return err
	}

	item.ID = invItem.ID.String()
	item.RecipeUnit = invItem.RecipeUnit
	item.RecipeUnitConversion = invItem.RecipeUnitConversion

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
			var vendorName string
			for k := range item.VendorItemMap {
				vendorName = k
				break
			}
			if vendorName != "" {
				item.SetSelectedVendor(vendorName)
			}
		} else {
			item.SetSelectedVendor(vendor.Name)
		}

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

	apiItem := NewItemAPI(*item)
	return &apiItem, nil
}

func (p *Presenter) GetFullInventoryItem(id string) (*ItemAPI, error) {
	item, err := p.GetInventoryItem(id)
	if err != nil {
		return nil, err
	}
	vendors, err := p.GetVendors()
	if err != nil {
		return nil, err
	}

	item.VendorItemMap = GetVendorMap(item.InventoryItemID, vendors)
	return item, nil
}

func (p *Presenter) UpdateFullInventoryItem(item *ItemAPI) (*validate.Errors, error) {
	verrs, err := p.UpdateBaseInventoryItem(item)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	err = p.deleteVendorItems(item)
	if err != nil {
		return verrs, err
	}

	return p.updateVendorItems(item)
}

func (p *Presenter) UpdateInventoryItem(item *ItemAPI) (*validate.Errors, error) {
	verrs, err := p.UpdateBaseInventoryItem(item)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return p.updateVendorItem(item)
}

func (p *Presenter) UpdateBaseInventoryItem(item *ItemAPI) (*validate.Errors, error) {
	invItem, err := ConvertToModelInventoryItem(item)
	if err != nil {
		return validate.NewErrors(), err
	}

	verrs, err := p.updateInventoryItemIndices(item)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return logic.UpdateInventoryItem(invItem, p.tx)
}

func (p *Presenter) updateInventoryItemIndices(item *ItemAPI) (*validate.Errors, error) {
	nextItems, err := logic.GetInvItemsAfter(item.ID, item.Index, p.tx)
	if err != nil {
		return validate.NewErrors(), err
	}

	for i, otherItem := range *nextItems {
		if otherItem.Index <= i+item.Index {
			otherItem.Index = i + item.Index + 1
			verrs, err := logic.UpdateBaseInventoryItem(&otherItem, p.tx)
			if verrs.HasAny() || err != nil {
				return verrs, err
			}
		}
	}

	return validate.NewErrors(), nil
}

func (p *Presenter) InsertInventoryItem(item *ItemAPI) (*validate.Errors, error) {
	invItem, err := ConvertToModelInventoryItem(item)
	if err != nil {
		return nil, err
	}

	verrs, err := logic.InsertInventoryItem(invItem, p.tx)
	if verrs.HasAny() && strings.Contains(verrs.Error(), "already exists") {
		verrs = validate.NewErrors()
		invItem, err = logic.ResurrectInventoryItem(invItem.Name, p.tx)
	}
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

func (p *Presenter) DestroyOrDeactivateInventoryItem(item *ItemAPI) error {
	inventoryItem, err := ConvertToModelInventoryItem(item)
	if err != nil {
		return err
	}
	if logic.HistoricalItemExists(item.ID, p.tx) {
		err = logic.DeactivateInventoryItem(inventoryItem, p.tx)
	} else {
		err = logic.DestroyInventoryItem(inventoryItem, p.tx)
	}
	if err != nil {
		return err
	}

	return p.destroyVendorMap(item.VendorItemMap)
}

func (p *Presenter) destroyVendorMap(vendorMap map[string]ItemAPI) error {
	for _, item := range vendorMap {
		vendorItem, err := ConvertToModelVendorItem(item, uuid.UUID{})
		if err != nil {
			return err
		}
		err = logic.DestroyVendorItem(vendorItem, p.tx)
		if err != nil {
			return err
		}
	}

	return nil
}

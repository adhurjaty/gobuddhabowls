package presentation

import (
	"buddhabowls/logic"
	"buddhabowls/models"
	"database/sql"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

func (p *Presenter) GetVendors() (*VendorsAPI, error) {
	vendors, err := logic.GetAllVendors(p.tx)
	if err != nil {
		return nil, err
	}

	apiVendors := NewVendorsAPI(vendors)
	return &apiVendors, nil
}

func (p *Presenter) GetVendor(id string) (*VendorAPI, error) {
	vendor, err := logic.GetVendor(id, p.tx)
	if err != nil {
		return nil, err
	}

	apiVendor := NewVendorAPI(vendor)
	return &apiVendor, nil
}

func (p *Presenter) InsertVendor(vendorAPI *VendorAPI) (*validate.Errors, error) {
	vendor, err := ConvertToModelVendor(vendorAPI)
	if err != nil {
		return nil, err
	}
	return logic.InsertVendor(vendor, p.tx)
}

func (p *Presenter) UpdateVendor(vendAPI *VendorAPI) (*validate.Errors, error) {
	vendor, err := ConvertToModelVendor(vendAPI)
	if err != nil {
		return nil, err
	}

	return logic.UpdateVendor(vendor, p.tx)
}

func (p *Presenter) UpdateVendorNoItems(vendorAPI *VendorAPI) (*validate.Errors, error) {
	vendor, err := ConvertToModelVendor(vendorAPI)
	if err != nil {
		return nil, err
	}

	return logic.UpdateVendorNoItems(vendor, p.tx)
}

func (p *Presenter) updateVendorItem(item *ItemAPI) (*validate.Errors, error) {
	vendorItem, err := p.getVendorItem(item)
	if err != nil {
		return nil, err
	}

	return logic.UpdateVendorItem(vendorItem, p.tx)
}

func (p *Presenter) updateVendorItems(item *ItemAPI) (*validate.Errors, error) {
	for vendorName, vItem := range item.VendorItemMap {
		vItem.InventoryItemID = item.ID
		vItem.SelectedVendor = vendorName

		dbVendorItem, err := logic.GetVendorItem(vItem.ID, p.tx)
		if err != nil {
			verrs, err := p.insertVendorItem(&vItem)
			if verrs.HasAny() || err != nil {
				return verrs, err
			}
		}

		vendorItem, err := ConvertToModelVendorItem(vItem, dbVendorItem.VendorID)

		verrs, err := logic.UpdateVendorItem(vendorItem, p.tx)
		if err != nil {
		}
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return validate.NewErrors(), nil
}

func (p *Presenter) insertVendorItem(item *ItemAPI) (*validate.Errors, error) {
	vendor, err := logic.GetVendorByName(item.SelectedVendor, p.tx)
	if err != nil {
		return validate.NewErrors(), err
	}

	vendorItem, err := ConvertToModelVendorItem(*item, vendor.ID)
	if err != nil {
		return validate.NewErrors(), err
	}

	return logic.InsertVendorItem(vendorItem, p.tx)
}

func (p *Presenter) getVendorItem(item *ItemAPI) (*models.VendorItem, error) {
	vendorID, err := p.getVendorID(item.SelectedVendor)
	if err != nil {
		return nil, err
	}
	vendorItem, err := ConvertToModelVendorItem(*item, vendorID)
	if err != nil {
		return nil, err
	}
	dbItem, err := logic.GetVendorItemByInvItem(item.InventoryItemID,
		vendorID.String(), p.tx)
	if err != nil {
		return nil, err
	}
	vendorItem.ID = dbItem.ID

	return vendorItem, nil
}

func (p *Presenter) getVendorID(name string) (uuid.UUID, error) {
	vendor, err := logic.GetVendorByName(name, p.tx)
	if err != nil {
		return uuid.UUID{}, err
	}

	return vendor.ID, nil
}

func (p *Presenter) GetBlankVendorItems() (*ItemsAPI, error) {
	vendors, err := logic.GetAllVendors(p.tx)
	if err != nil {
		return nil, err
	}

	items := ItemsAPI{}
	for _, vendor := range *vendors {
		vItem := models.VendorItem{
			VendorID:   vendor.ID,
			Vendor:     vendor,
			Conversion: 1,
			PurchasedUnit: sql.NullString{
				String: "EA",
				Valid:  true,
			},
		}
		item := NewItemAPI(vItem)
		item.SelectedVendor = vendor.Name
		items = append(items, item)
	}

	return &items, err
}

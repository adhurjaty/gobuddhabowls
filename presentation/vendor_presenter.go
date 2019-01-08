package presentation

import (
	"buddhabowls/logic"
	"buddhabowls/models"
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

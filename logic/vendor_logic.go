package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"sort"
)

func GetAllVendors(tx *pop.Connection) (*models.Vendors, error) {
	factory := models.ModelFactory{}
	vendors := &models.Vendors{}
	err := factory.CreateModelSlice(vendors, tx.Eager().Q())
	if err != nil {
		return nil, err
	}

	sort.Slice(*vendors, func(i, j int) bool {
		return (*vendors)[i].Name < (*vendors)[j].Name
	})

	return vendors, err
}

func GetVendor(id string, tx *pop.Connection) (*models.Vendor, error) {
	factory := models.ModelFactory{}
	vendor := &models.Vendor{}
	err := factory.CreateModel(vendor, tx, id)

	return vendor, err
}

func GetLatestVendor(invItemId string, tx *pop.Connection) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	query := tx.Where("vi.inventory_item_id = ?", invItemId).
		Join("vendor_items vi", "vendors.id=vi.vendor_id").
		Join("order_items oi", "vi.inventory_item_id=oi.inventory_item_id").
		Join("purchase_orders po", "po.id=oi.order_id").
		Order("po.order_date desc")

	err := query.First(vendor)

	return vendor, err
}

func InsertVendor(vendor *models.Vendor, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := tx.ValidateAndCreate(vendor)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	// insert items
	for _, item := range vendor.Items {
		// need to ensure that items don't get the vendor item ID
		item.ID = uuid.UUID{}
		item.VendorID = vendor.ID
		verrs, err = tx.ValidateAndCreate(&item)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return verrs, nil
}

func UpdateVendor(vendor *models.Vendor, tx *pop.Connection) (*validate.Errors, error) {
	oldVendor, err := GetVendor(vendor.ID.String(), tx)
	if err != nil {
		return nil, err
	}

	verrs, err := tx.ValidateAndUpdate(vendor)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	oldItems := oldVendor.Items
	containsFunc := func(item models.VendorItem, itemArr models.VendorItems) bool {
		for _, otherItem := range itemArr {
			if item.ID == otherItem.ID {
				return true
			}
		}
		return false
	}

	// update or insert items
	for _, item := range vendor.Items {
		item.VendorID = vendor.ID
		if containsFunc(item, oldItems) {
			verrs, err = tx.ValidateAndUpdate(&item)
		} else {
			verrs, err = tx.ValidateAndCreate(&item)
		}
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	// delete items are removed from the order list
	for _, item := range oldItems {
		if !containsFunc(item, vendor.Items) {
			err = tx.Destroy(&item)
			if err != nil {
				return verrs, err
			}
		}
	}

	return verrs, nil
}

func UpdateVendorItems(items *models.VendorItems, tx *pop.Connection) (*validate.Errors, error) {
	for _, item := range *items {
		verrs, err := tx.ValidateAndUpdate(&item)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return validate.NewErrors(), nil
}

func DeleteVendor(vendor *models.Vendor, tx *pop.Connection) error {
	for _, item := range vendor.Items {
		tx.Destroy(&item)
	}
	return tx.Destroy(vendor)
}

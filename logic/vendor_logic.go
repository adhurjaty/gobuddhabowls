package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
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

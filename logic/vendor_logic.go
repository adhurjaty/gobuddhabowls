package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
)

func GetVendor(id string, tx *pop.Connection) (*models.Vendor, error) {
	factory := models.ModelFactory{}
	vendor := &models.Vendor{}
	err := factory.CreateModel(vendor, tx, id)

	return vendor, err
}

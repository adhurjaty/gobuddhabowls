package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

func GetInventoryItems(tx *pop.Connection) (*models.InventoryItems, error) {
	factory := models.ModelFactory{}
	items := &models.InventoryItems{}
	err := factory.CreateModelSlice(items, tx.Eager().Q())

	if err != nil {
		return nil, err
	}

	items.Sort()

	return items, err
}

func GetInventoryItem(id string, tx *pop.Connection) (*models.InventoryItem, error) {
	factory := models.ModelFactory{}
	item := &models.InventoryItem{}
	err := factory.CreateModel(item, tx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func UpdateInventoryItem(item *models.InventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndUpdate(item)
}

func InsertInventoryItem(item *models.InventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := updateIndices(item.Index, tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return tx.ValidateAndCreate(item)
}

func updateIndices(index int, tx *pop.Connection) (*validate.Errors, error) {
	items, err := GetInventoryItems(tx)
	if err != nil {
		return validate.NewErrors(), err
	}

	for _, item := range *items {
		if item.Index < index {
			continue
		}

		item.Index++
		verrs, err := tx.ValidateAndUpdate(&item)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return validate.NewErrors(), nil
}

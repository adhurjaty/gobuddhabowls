package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
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

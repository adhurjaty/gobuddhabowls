package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"sort"
)

func GetAllInventories(tx *pop.Connection) (*models.Inventories, error) {
	factory := models.ModelFactory{}
	inventories := &models.Inventories{}
	err := factory.CreateModelSlice(inventories, tx.Eager().Q())
	if err != nil {
		return nil, err
	}

	sort.Slice(*inventories, func(i, j int) bool {
		return (*inventories)[i].Date.Unix() > (*inventories)[j].Date.Unix()
	})

	return inventories, err
}

func GetInventory(id string, tx *pop.Connection) (*models.Inventory, error) {
	factory := models.ModelFactory{}
	inventory := &models.Inventory{}
	err := factory.CreateModel(inventory, tx, id)

	return inventory, err
}

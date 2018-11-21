package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"time"
)

func GetAllInventories(tx *pop.Connection) (*models.Inventories, error) {
	factory := models.ModelFactory{}
	inventories := &models.Inventories{}
	err := factory.CreateModelSlice(inventories, tx.Eager().Order("date DESC"))
	if err != nil {
		return nil, err
	}

	return inventories, err
}

func GetInventories(startTime time.Time, endTime time.Time, tx *pop.Connection) (*models.Inventories, error) {
	factory := models.ModelFactory{}
	inventories := &models.Inventories{}

	startVal := startTime.Format(time.RFC3339)
	endVal := endTime.Format(time.RFC3339)
	query := tx.Eager().Where("date >= ? AND date < ?",
		startVal, endVal).Order("date DESC")

	err := factory.CreateModelSlice(inventories, query)
	if err != nil {
		return nil, err
	}

	return inventories, err
}

func GetInventory(id string, tx *pop.Connection) (*models.Inventory, error) {
	factory := models.ModelFactory{}
	inventory := &models.Inventory{}
	err := factory.CreateModel(inventory, tx, id)

	return inventory, err
}

package logic

import (
	"buddhabowls/models"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
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

func GetLatestInventory(date time.Time, tx *pop.Connection) (*models.Inventory, error) {
	factory := models.ModelFactory{}
	inventories := &models.Inventories{}
	query := tx.Eager().Where("date <= ?",
		date.Format(time.RFC3339)).Order("date DESC").Limit(1)
	err := factory.CreateModelSlice(inventories, query)
	if err != nil {
		return nil, err
	}
	if len(*inventories) == 0 {
		return nil, nil
	}
	inventory := (*inventories)[0]
	return &inventory, nil
}

func GetCountInventoryItemByInvItem(invItemID string, inventoryID string, tx *pop.Connection) (*models.CountInventoryItem, error) {
	item := &models.CountInventoryItem{}
	query := tx.Eager().Where("inventory_item_id = ?", invItemID).
		Where("inventory_id = ?", inventoryID)
	err := query.First(item)
	return item, err
}

func InsertInventory(inventory *models.Inventory, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := tx.ValidateAndCreate(inventory)
	if err != nil {
		return verrs, err
	}

	for _, item := range inventory.Items {
		item.InventoryID = inventory.ID
		item.ID, err = uuid.NewV4()
		if err != nil {
			return nil, err
		}
		verrs, err = tx.ValidateAndCreate(&item)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}
	for _, item := range inventory.PrepItems {
		item.InventoryID = inventory.ID
		item.ID = uuid.UUID{}
		verrs, err = tx.ValidateAndCreate(&item)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return verrs, nil
}

func UpdateInventory(inventory *models.Inventory, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := tx.ValidateAndUpdate(inventory)
	if err != nil {
		return verrs, err
	}

	for _, item := range inventory.Items {
		item.InventoryID = inventory.ID
		verrs, err = tx.ValidateAndUpdate(&item)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return verrs, nil
}

func UpdateCountInventoryItem(item *models.CountInventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndUpdate(item)
}

func DestroyInventory(inventory *models.Inventory, tx *pop.Connection) error {
	for _, item := range inventory.Items {
		tx.Destroy(&item)
	}
	return tx.Destroy(inventory)
}

func InsertCountInventoryItem(item *models.CountInventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(item)
}

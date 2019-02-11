package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

func GetInventoryItems(tx *pop.Connection) (*models.InventoryItems, error) {
	factory := models.ModelFactory{}
	items := &models.InventoryItems{}
	q := tx.Eager().Where("is_active = true")
	err := factory.CreateModelSlice(items, q)

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

func GetInvItemsAfter(id string, idx int, tx *pop.Connection) (*models.InventoryItems, error) {
	items := &models.InventoryItems{}
	query := tx.Where("id != ?", id).Where("index >= ?", idx).
		Order("index")
	err := query.All(items)

	return items, err
}

func HistoricalItemExists(inventoryItemID string, tx *pop.Connection) bool {
	orderItem := &models.OrderItem{}
	countItem := &models.CountInventoryItem{}

	q := tx.Where("inventory_item_id = ?", inventoryItemID)

	return q.First(orderItem) == nil ||
		q.First(countItem) == nil
}

func UpdateInventoryItem(item *models.InventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := updateIndices(item, tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}
	return tx.ValidateAndUpdate(item)
}

func UpdateBaseInventoryItem(item *models.InventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndUpdate(item)
}

func InsertInventoryItem(item *models.InventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := updateIndices(item, tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return tx.ValidateAndCreate(item)
}

func updateIndices(invItem *models.InventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	items, err := GetInventoryItems(tx)
	if err != nil {
		return validate.NewErrors(), err
	}

	offset := 0
	for i, item := range *items {
		if item.ID.String() == invItem.ID.String() {
			offset = -1
			continue
		}
		if item.Index == invItem.Index {
			offset = 1
		}

		item.Index = i + offset
		verrs, err := tx.ValidateAndUpdate(&item)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return validate.NewErrors(), nil
}

func DeactivateInventoryItem(item *models.InventoryItem, tx *pop.Connection) error {
	item.IsActive = false
	_, err := tx.ValidateAndUpdate(item)
	return err
}

func DestroyInventoryItem(item *models.InventoryItem, tx *pop.Connection) error {
	return tx.Destroy(item)
}

func ResurrectInventoryItem(name string, tx *pop.Connection) (*models.InventoryItem, error) {
	item := &models.InventoryItem{}
	err := tx.Where("name = ?", name).First(item)
	if err != nil {
		return nil, err
	}
	item.IsActive = true
	_, err = tx.ValidateAndUpdate(item)
	return item, err
}

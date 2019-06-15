package logic

import (
	"buddhabowls/models"
	"errors"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

func GetInventoryItems(tx *pop.Connection) (*models.InventoryItems, error) {
	return getInventoryItemsHelper(tx.Q())
}

func GetInventoryItemsOfCategory(id string, catID string, tx *pop.Connection) (*models.InventoryItems, error) {
	query := tx.Where("category_id = ?", catID).
		Where("id != ?", id)
	return getInventoryItemsHelper(query)
}

func getInventoryItemsHelper(query *pop.Query) (*models.InventoryItems, error) {
	query = query.Where("is_active = true")
	factory := models.ModelFactory{}
	items := &models.InventoryItems{}
	err := factory.CreateModelSlice(items, query)

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

func HistoricalItemExists(inventoryItemID string, tx *pop.Connection) bool {
	orderItem := &models.OrderItem{}
	countItem := &models.CountInventoryItem{}

	q := tx.Where("inventory_item_id = ?", inventoryItemID)

	return q.First(orderItem) == nil ||
		q.First(countItem) == nil
}

func UpdateInventoryItem(item *models.InventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := UpdateIndices(item, tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}
	return tx.ValidateAndUpdate(item)
}

func UpdateBaseItem(item models.GenericItem, tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndUpdate(item)
}

func InsertInventoryItem(item *models.InventoryItem, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := UpdateIndices(item, tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return tx.ValidateAndCreate(item)
}

func UpdateIndices(genItem models.GenericItem, tx *pop.Connection) (*validate.Errors, error) {
	var items models.GenericItems
	var err error
	verrs := validate.NewErrors()

	switch genItem.(type) {
	case *models.InventoryItem:
		items, err = GetInventoryItemsOfCategory(genItem.GetID().String(),
			genItem.GetCategoryID().String(), tx)
		if err != nil {
			return verrs, err
		}
		break
	case *models.PrepItem:
		prepItem := genItem.(*models.PrepItem)
		items, err = GetPrepItemsOfCategory(prepItem, tx)
		if err != nil {
			return verrs, err
		}
		break
	default:
		return verrs, errors.New("unimplemented type")
	}

	offset := 0
	for i, item := range *items.ToGenericItems() {
		if i == genItem.GetIndex() {
			offset = 1
		}

		item.SetIndex(i + offset)
		verrs, err := tx.ValidateAndUpdate(item)
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

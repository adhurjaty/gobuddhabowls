package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"time"
)

// GetPurchaseOrders retrieves purchase orders from within the given start and end dates
func GetPurchaseOrders(startTime, endTime time.Time, tx *pop.Connection) (*models.PurchaseOrders, error) {
	startVal := startTime.Format(time.RFC3339)
	endVal := endTime.Format(time.RFC3339)

	q := tx.Eager().Where("order_date >= ? AND order_date <= ?",
		startVal, endVal).Order("order_date DESC")

	pos := &models.PurchaseOrders{}
	factory := models.ModelFactory{}

	err := factory.CreateModelSlice(pos, q)

	return pos, err
}

// GetPurchaseOrder returns a PurchaseOrder model
func GetPurchaseOrder(id string, tx *pop.Connection) (*models.PurchaseOrder, error) {
	factory := models.ModelFactory{}
	po := &models.PurchaseOrder{}
	err := factory.CreateModel(po, tx, id)

	return po, err
}

func InsertPurchaseOrder(purchaseOrder *models.PurchaseOrder, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := tx.ValidateAndCreate(purchaseOrder)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	// insert items
	for _, item := range purchaseOrder.Items {
		item.OrderID = purchaseOrder.ID
		verrs, err = tx.ValidateAndCreate(&item)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return verrs, nil
}

func UpdatePurchaseOrder(purchaseOrder *models.PurchaseOrder, tx *pop.Connection) (*validate.Errors, error) {
	oldPO, err := GetPurchaseOrder(purchaseOrder.ID.String(), tx)
	if err != nil {
		return nil, err
	}

	verrs, err := tx.ValidateAndUpdate(purchaseOrder)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	oldItems := oldPO.Items
	containsFunc := func(item models.OrderItem, itemArr models.OrderItems) bool {
		for _, otherItem := range itemArr {
			if item.ID == otherItem.ID {
				return true
			}
		}
		return false
	}

	// update or insert items
	for _, item := range purchaseOrder.Items {
		item.OrderID = purchaseOrder.ID
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
		if !containsFunc(item, purchaseOrder.Items) {
			err = tx.Destroy(item)
			if err != nil {
				return verrs, err
			}
		}
	}

	return verrs, nil
}

func DeletePurchaseOrder(purchaseOrder *models.PurchaseOrder, tx *pop.Connection) error {
	return tx.Destroy(purchaseOrder)
}

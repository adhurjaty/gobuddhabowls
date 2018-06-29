package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"time"
)

// GetPurchaseOrders retrieves purchase orders from within the given start and end dates
func GetPurchaseOrders(startTime, endTime time.Time, tx *pop.Connection) (*models.PurchaseOrders, error) {
	startTime = OffsetStart(startTime)
	endTime = OffsetEnd(endTime)

	startVal := startTime.Format(time.RFC3339)
	endVal := endTime.Format(time.RFC3339)

	q := tx.Eager().Where("order_date >= ? AND order_date < ?",
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
	return tx.Eager().ValidateAndCreate(purchaseOrder)
}

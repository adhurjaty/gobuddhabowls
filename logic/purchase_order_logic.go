package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"time"
)

// GetPurchaseOrders retrieves purchase orders from within the given start and end dates
func GetPurchaseOrders(startTime, endTime time.Time, tx *pop.Connection) (*models.PurchaseOrders, error) {
	startTime = UnoffsetStart(startTime)
	endTime = UnoffsetEnd(endTime)

	startVal := startTime.Format(time.RFC3339)
	endVal := endTime.Format(time.RFC3339)

	q := tx.Eager().Where("order_date >= ? AND order_date < ?",
		startVal, endVal).Order("order_date DESC")

	pos := &models.PurchaseOrders{}
	factory := models.ModelFactory{}

	err := factory.CreateModelSlice(pos, q)

	return pos, err
}

package presentationlayer

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
)

// GetOpenRecPurchaseOrders gets the purchase orders give a query and returns
// them separated by whether they were received
func GetOpenRecPurchaseOrders(q *pop.Query) (models.PurchaseOrders, models.PurchaseOrders, error) {
	purchaseOrders, err := models.LoadPurchaseOrders(q)
	if err != nil {
		return nil, nil, err
	}

	var openPos models.PurchaseOrders
	var recPos models.PurchaseOrders

	for _, po := range *purchaseOrders {
		if po.ReceivedDate.Valid {
			recPos = append(recPos, po)
		} else {
			openPos = append(openPos, po)
		}
	}

	return openPos, recPos, nil
}

// func GetCategoryDetailsAndTotal(open *models.PurchaseOrders, rec *models.PurchaseOrders) (){

// }

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

// GetOrderCategoryDetails gets a CategoryBreakdown from a purchase order
// func GetOrderCategoryDetails(po models.PurchaseOrder) models.CategoryBreakdown {
// 	breakdownMap := po.GetCategoryCosts()

// 	return fromCategoryMap(breakdownMap)
// }

// GetAllCategoryDetails gets a category breakdown of all orders
// may expand to return open, rec and total
func GetAllCategoryDetails(open models.PurchaseOrders, rec models.PurchaseOrders) models.CategoryBreakdown {
	returnBreakdown := models.CategoryBreakdown{}

	// for now just combine the purchase orders
	for _, po := range open {
		returnBreakdown.Add(po.GetCategoryCosts())
	}
	for _, po := range rec {
		returnBreakdown.Add(po.GetCategoryCosts())
	}

	return returnBreakdown
}

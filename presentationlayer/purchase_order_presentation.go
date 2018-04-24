package presentationlayer

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
)

// CategoryBreakdown is used to display category values in the
// summary table on the purchase orders index page
type CategoryBreakdown struct {
	Total      float64
	Categories []CategoryBreakdownItem
}

// CategoryBreakdownItem is an item with the breakdown
type CategoryBreakdownItem struct {
	Value    float64
	Category models.InventoryItemCategory
}

// Add adds CategoryBreakdown by combining the totals for each category
func (cb *CategoryBreakdown) Add(otherCb CategoryBreakdown) {
	combinedMaps := models.CombineCategoryMaps(cb.toCategoryMap(), otherCb.toCategoryMap())
	newCb := fromCategoryMap(combinedMaps)
	*cb = newCb
}

func (cb CategoryBreakdown) toCategoryMap() map[models.InventoryItemCategory]float64 {
	cbMap := make(map[models.InventoryItemCategory]float64)
	for _, item := range cb.Categories {
		cbMap = models.AddToCategoryMap(cbMap, item.Category, item.Value)
	}

	return cbMap
}

func fromCategoryMap(m map[models.InventoryItemCategory]float64) CategoryBreakdown {
	returnBreakdown := CategoryBreakdown{Categories: []CategoryBreakdownItem{}}
	var total float64
	for cat, value := range m {
		returnBreakdown.Categories = append(returnBreakdown.Categories,
			CategoryBreakdownItem{
				Category: cat,
				Value:    value,
			})
		total += value
	}
	returnBreakdown.Total = total

	return returnBreakdown
}

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
func GetOrderCategoryDetails(po models.PurchaseOrder) CategoryBreakdown {
	breakdownMap := po.GetCategoryCosts()

	return fromCategoryMap(breakdownMap)
}

// GetAllCategoryDetails gets a category breakdown of all orders
// may expand to return open, rec and total
func GetAllCategoryDetails(open models.PurchaseOrders, rec models.PurchaseOrders) CategoryBreakdown {
	returnBreakdown := CategoryBreakdown{}

	// for now just combine the purchase orders
	for _, po := range open {
		returnBreakdown.Add(fromCategoryMap(po.GetCategoryCosts()))
	}
	for _, po := range rec {
		returnBreakdown.Add(fromCategoryMap(po.GetCategoryCosts()))
	}

	return returnBreakdown
}

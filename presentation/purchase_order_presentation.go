package presentation

import (
	"buddhabowls/models"
	"fmt"
	"github.com/gobuffalo/pop"
	"sort"
	"strings"
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

// GetBarChartJSONData gets a JSON string for showing bar chart category breakdown
func GetBarChartJSONData(open models.PurchaseOrders, rec models.PurchaseOrders) string {
	breakdown := GetAllCategoryDetails(open, rec)
	jsonItems := make([]string, len(breakdown.Categories))

	for i, item := range breakdown.Categories {
		jsonItems[i] = fmt.Sprintf("{\"Name\":\"%s\",\"Value\":%f,\"Background\":\"%s\"}",
			item.Category.Name, item.Value, item.Category.Background)
	}

	return "[" + strings.Join(jsonItems, ",") + "]"
}

// GetLineChartJSONData gets a JSON string for showing the line chart category breakdown
func GetLineChartJSONData(open models.PurchaseOrders, rec models.PurchaseOrders) string {
	var jsonItems []string
	categoryBreakdownCache := make([]models.CategoryBreakdown, len(open)+len(rec))

	combinedSortedPos := append(open, rec...)
	sort.Slice(combinedSortedPos, func(i, j int) bool {
		return combinedSortedPos[i].OrderDate.Time.Unix() < combinedSortedPos[j].OrderDate.Time.Unix()
	})

	// only show categories that are in the provided purchase orders
	// show 0 category value for unused ones in particular po's
	categoriesMap := make(map[models.InventoryItemCategory]bool)
	for i, po := range combinedSortedPos {
		breakdown := po.GetCategoryCosts()
		categoryBreakdownCache[i] = breakdown
		for _, item := range breakdown.Categories {
			categoriesMap[item.Category] = true
		}
	}

	// extract and sort categories
	var categories models.InventoryItemCategories
	for category := range categoriesMap {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Index < categories[j].Index
	})

	for i, po := range combinedSortedPos {
		breakdownMap := categoryBreakdownCache[i].ToCategoryMap()
		for _, category := range categories {

			value, ok := breakdownMap[category]
			if !ok {
				value = 0
			}

			jsonItems = append(jsonItems,
				fmt.Sprintf("{\"Name\":\"%s\",\"Date\":\"%s\",\"Value\":%f,\"Background\":\"%s\"}",
					category.Name, po.OrderDate.Time.UTC(), value, category.Background))
		}
	}

	return "[" + strings.Join(jsonItems, ",") + "]"
}

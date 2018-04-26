package models

import (
	"sort"
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
	Category InventoryItemCategory
}

// Add adds CategoryBreakdown by combining the totals for each category
func (cb *CategoryBreakdown) Add(otherCb CategoryBreakdown) {
	combinedMaps := CombineCategoryMaps(cb.toCategoryMap(), otherCb.toCategoryMap())
	newCb := FromCategoryMap(combinedMaps)
	*cb = newCb
}

func (cb CategoryBreakdown) toCategoryMap() map[InventoryItemCategory]float64 {
	cbMap := make(map[InventoryItemCategory]float64)
	for _, item := range cb.Categories {
		cbMap = AddToCategoryMap(cbMap, item.Category, item.Value)
	}

	return cbMap
}

// FromCategoryMap converts a map of categories to aggregate value into a CategoryBreakdown
func FromCategoryMap(m map[InventoryItemCategory]float64) CategoryBreakdown {
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

	sort.Slice(returnBreakdown.Categories, func(i, j int) bool {
		return returnBreakdown.Categories[i].Category.Index < returnBreakdown.Categories[j].Category.Index
	})

	return returnBreakdown
}

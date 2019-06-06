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
	Category ItemCategory
}

// Add adds CategoryBreakdown by combining the totals for each category
func (cb *CategoryBreakdown) Add(otherCb CategoryBreakdown) {
	combinedMaps := CombineCategoryMaps(cb.ToCategoryMap(), otherCb.ToCategoryMap())
	newCb := FromCategoryMap(combinedMaps)
	*cb = newCb
}

// ToCategoryMap converts the breakdown to a map of category to value
func (cb CategoryBreakdown) ToCategoryMap() map[ItemCategory]float64 {
	cbMap := make(map[ItemCategory]float64)
	for _, item := range cb.Categories {
		cbMap = AddToCategoryMap(cbMap, item.Category, item.Value)
	}

	return cbMap
}

// FromCategoryMap converts a map of categories to aggregate value into a CategoryBreakdown
func FromCategoryMap(m map[ItemCategory]float64) CategoryBreakdown {
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
		return returnBreakdown.Categories[i].Category.GetIndex() <
			returnBreakdown.Categories[j].Category.GetIndex()
	})

	return returnBreakdown
}

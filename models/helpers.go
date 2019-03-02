package models

import (
	"buddhabowls/helpers"
	"database/sql"
	"github.com/gobuffalo/pop"
)

// AddToCategoryMap adds to map the key and value and returns map
func AddToCategoryMap(m map[ItemCategory]float64, key ItemCategory, value float64) map[ItemCategory]float64 {
	if v, ok := m[key]; ok {
		m[key] = v + value
	} else {
		m[key] = value
	}
	return m
}

// CombineCategoryMaps adds 2 maps together by combining values for each key
func CombineCategoryMaps(m1 map[ItemCategory]float64, m2 map[ItemCategory]float64) map[ItemCategory]float64 {
	for k, v := range m2 {
		m1 = AddToCategoryMap(m1, k, v)
	}

	return m1
}

// GetYears gets the years for which there is company data
func GetYears(tx *pop.Connection) ([]int, error) {
	yearResult := make([]int, 50)
	// Search for just the years in purchase orders
	q := tx.RawQuery("SELECT DISTINCT EXTRACT(YEAR FROM order_date) FROM purchase_orders ORDER BY EXTRACT(YEAR FROM order_date) ASC")

	if err := q.All(&yearResult); err != nil {
		return nil, err
	}

	// throw away extra allocated data. Probably a better way to do this
	years := []int{}
	for _, val := range yearResult {
		if val > 2000 {
			years = append(years, val)
		}
	}
	thisYear := helpers.Today().Year()
	if len(years) > 0 && years[len(years)-1] < thisYear {
		years = append(years, thisYear)
	}

	return years, nil
}

func StringToNullString(s string) sql.NullString {
	valid := len(s) > 0
	return sql.NullString{
		Valid:  valid,
		String: s,
	}
}

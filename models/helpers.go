package models

// AddToCategoryMap adds to map the key and value and returns map
func AddToCategoryMap(m map[InventoryItemCategory]float64, key InventoryItemCategory, value float64) map[InventoryItemCategory]float64 {
	if v, ok := m[key]; ok {
		m[key] = v + value
	} else {
		m[key] = value
	}
	return m
}

// CombineCategoryMaps adds 2 maps together by combining values for each key
func CombineCategoryMaps(m1 map[InventoryItemCategory]float64, m2 map[InventoryItemCategory]float64) map[InventoryItemCategory]float64 {
	for k, v := range m2 {
		m1 = AddToCategoryMap(m1, k, v)
	}

	return m1
}

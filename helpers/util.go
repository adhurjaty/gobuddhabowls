package helpers

func AddToMap(m map[string]float64, key string, value float64) map[string]float64 {
	if v, ok := m[key]; ok {
		m[key] = v + value
	} else {
		m[key] = value
	}
	return m
}

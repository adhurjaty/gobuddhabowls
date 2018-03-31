package helpers

import (
	"fmt"
	"time"
)

// AddToMap adds a float value to the map at key, if no key is present, create a
// new entry with value as the value
func AddToMap(m map[string]float64, key string, value float64) map[string]float64 {
	if v, ok := m[key]; ok {
		m[key] = v + value
	} else {
		m[key] = value
	}
	return m
}

// FormatDate returns the date supplied into mm/dd/yyyy format
func FormatDate(date time.Time) string {
	year, month, day := date.Date()
	return fmt.Sprintf("%02d/%02d/%d", month, day, year)
}

package helpers

import (
	"fmt"
	"github.com/gobuffalo/uuid"
	"time"
)

// AddToMap adds a float value to the map at key, if no key is present, create a
// new entry with value as the value
func AddToMap(m map[interface{}]float64, key interface{}, value float64) map[interface{}]float64 {
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

// RFC3339Date converts date to a parseable format
func RFC3339Date(date time.Time) string {
	return date.Format(time.RFC3339)
}

func Today() time.Time {
	nowTime := time.Now()
	return time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(),
		0, 0, 0, 0, nowTime.Location())
}

func IsBlankUUID(id uuid.UUID) bool {
	return id.String() == uuid.UUID{}.String()
}

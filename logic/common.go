package logic

import (
	"time"
)

const dayStart = 4 * time.Hour

// OffsetStart turns a start of a day into a business-specific start of the day
// Dan's shop closes at 2 AM and wants sales to go into the previous day
func OffsetStart(t time.Time) time.Time {
	return t.Add(dayStart)
}

// OffsetEnd turns a start of a day into a business-specific end of the day
func OffsetEnd(t time.Time) time.Time {
	return t.Add(dayStart - time.Nanosecond)
}

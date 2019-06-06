package logic

import (
	"errors"
	"fmt"
	"time"
)

const (
	weekStart = 1 // denotes that the first day of the week is Monday (not Sunday)
)

// PeriodSelector holds values for periods of the year
type PeriodSelector struct {
	Periods []Period
	Year    int
}

// NewPeriodSelector creates a new period selector given a year
func NewPeriodSelector(year int) PeriodSelector {
	p := PeriodSelector{}
	p.Year = year
	theFirst := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	firstWeekStartDiff := ((7 + weekStart) - int(theFirst.Weekday())) % 7

	p.Periods = []Period{}
	// startTime := theFirst.Add(dayStart)
	startTime := theFirst
	endTime := startTime
	index := 1

	if firstWeekStartDiff != 0 {
		endTime = startTime.AddDate(0, 0, firstWeekStartDiff)
		period := NewPeriod(startTime, endTime)
		period.Index = index
		index++
		p.Periods = append(p.Periods, period)
		startTime = endTime
	}

	for index < 15 {
		endTime = startTime.AddDate(0, 0, 28)

		// if the period goes to the end of the year
		if endTime.YearDay() < startTime.YearDay() {
			endTime = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		}
		period := NewPeriod(startTime, endTime)
		period.Index = index
		p.Periods = append(p.Periods, period)
		startTime = endTime
		index++
	}

	return p
}

// GetPeriod gets the period that contains the date provided
func (p PeriodSelector) GetPeriod(date time.Time) Period {
	for _, period := range p.Periods {
		if date.Unix() >= period.StartTime().Unix() && date.Unix() < period.EndTime().Unix() {
			return period
		}
	}

	// this should never be reached
	panic(errors.New("No period matching date"))
}

// GetWeek gest the week that contains the date provided
func (p PeriodSelector) GetWeek(date time.Time) Week {
	period := p.GetPeriod(date)
	return period.GetWeek(date)
}

// Period represents the 4 week business period (14 per year)
type Period struct {
	Index int
	Weeks []Week
}

// NewPeriod creates a new period
func NewPeriod(startTime time.Time, endTime time.Time) Period {
	p := Period{}
	dayDiff := int(endTime.Sub(startTime).Hours() / 24)
	fullWeeks := int(dayDiff / 7)
	p.Weeks = make([]Week, fullWeeks)

	weekEnd := startTime
	for i := 0; i < fullWeeks; i++ {
		weekEnd = startTime.AddDate(0, 0, 7)
		week := Week{StartTime: startTime, EndTime: weekEnd.Add(-time.Nanosecond), Index: i + 1}
		p.Weeks[i] = week
		startTime = weekEnd
	}

	if dayDiff%7 > 0 {
		week := Week{StartTime: startTime, EndTime: endTime, Index: fullWeeks + 1}
		p.Weeks = append(p.Weeks, week)
	}

	return p
}

// GetWeek get the week that contains the date
func (p Period) GetWeek(date time.Time) Week {
	for _, week := range p.Weeks {
		if date.Unix() >= week.StartTime.Unix() && date.Unix() < week.EndTime.Unix() {
			return week
		}
	}

	panic(errors.New("No week matching date"))
}

// StartTime gets the time at the beginning of the period
func (p Period) StartTime() time.Time {
	return p.Weeks[0].StartTime
}

// EndTime get the time at the end of the period
func (p Period) EndTime() time.Time {
	return p.Weeks[len(p.Weeks)-1].EndTime
}

func (p Period) String() string {
	start := p.StartTime()
	end := p.EndTime()
	return fmt.Sprintf("P%d %d/%d-%d/%d", p.Index, start.Month(), start.Day(),
		end.Month(), end.Day())
}

func (p Period) SelectValue() interface{} {
	if p.Weeks == nil {
		return ""
	}
	return p.StartTime().Format(time.RFC3339)
}

func (p Period) SelectLabel() string {
	if p.Weeks == nil {
		return ""
	}
	return p.String()
}

// Week represents a business week (Monday to Sunday)
type Week struct {
	Index     int
	StartTime time.Time
	EndTime   time.Time
}

func (w Week) String() string {
	start := w.StartTime
	end := w.EndTime
	return fmt.Sprintf("WK%d %d/%d-%d/%d", w.Index, start.Month(), start.Day(),
		end.Month(), end.Day())
}

func (w Week) SelectValue() interface{} {
	return w.StartTime.Format(time.RFC3339)
}

func (w Week) SelectLabel() string {
	if w.Index == 0 {
		return ""
	}
	return w.String()
}

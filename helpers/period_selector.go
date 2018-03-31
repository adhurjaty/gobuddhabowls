package helpers

import (
	"errors"
	"fmt"
	"time"
)

const weekStart = 1
const dayStart = time.Hour * 4

// PeriodSelector holds values for periods of the year
type PeriodSelector struct {
	Periods []Period
	Year    int
}

// Init populates a new period selector given a year
func (p *PeriodSelector) Init(year int) {
	p.Year = year
	theFirst := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	firstWeekStartDiff := ((7 + weekStart) - int(theFirst.Weekday())) % 7
	fmt.Println(firstWeekStartDiff)

	p.Periods = []Period{}
	startTime := theFirst.Add(dayStart)
	endTime := startTime
	index := 1
	var period Period
	if firstWeekStartDiff != 0 {
		period := *new(Period)
		endTime = startTime.AddDate(0, 0, firstWeekStartDiff)
		period.Init(startTime, endTime)
		period.Index = index
		index++
		p.Periods = append(p.Periods, period)
		startTime = endTime
	}

	for i := 0; i < 13; i++ {
		period = *new(Period)
		endTime = startTime.AddDate(0, 0, 28)

		// if the period goes to the end of the year
		if endTime.YearDay() < startTime.YearDay() {
			endTime = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(dayStart)
		}
		period.Init(startTime, endTime)
		period.Index = index
		p.Periods = append(p.Periods, period)
		startTime = endTime
		index++
	}
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

// Period represents the 4 week business period (13 per year)
type Period struct {
	Index int
	Weeks []Week
}

// Init populates a period
func (p *Period) Init(startTime time.Time, endTime time.Time) {
	dayDiff := int(endTime.Sub(startTime).Hours() / 24)
	fullWeeks := int(dayDiff / 7)
	p.Weeks = []Week{}

	var week Week
	weekEnd := startTime
	for i := 0; i < fullWeeks; i++ {
		weekEnd = startTime.AddDate(0, 0, 7)
		week = Week{StartTime: startTime, EndTime: weekEnd, Index: i + 1}
		p.Weeks = append(p.Weeks, week)
		startTime = weekEnd
	}

	if dayDiff%7 > 0 {
		week = Week{StartTime: startTime, EndTime: endTime, Index: fullWeeks + 1}
		p.Weeks = append(p.Weeks, week)
	}
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
	start := p.UnoffsetStart()
	end := p.UnoffsetEnd()
	return fmt.Sprintf("P%d %d/%d-%d/%d", p.Index, start.Month(), start.Day(),
		end.Month(), end.Day())
}

func (p Period) StartDateStr() string {
	return FormatDate(p.UnoffsetStart())
}

func (p Period) EndDateStr() string {
	return FormatDate(p.UnoffsetEnd())
}

func (p Period) UnoffsetStart() time.Time {
	return p.StartTime().Add(-dayStart)
}

func (p Period) UnoffsetEnd() time.Time {
	return p.EndTime().Add(-dayStart - time.Nanosecond)
}

func (p Period) Equals(other Period) bool {
	return p.StartTime().Unix() == other.StartTime().Unix()
}

// Week represents a business week (Monday to Sunday)
type Week struct {
	Index     int
	StartTime time.Time
	EndTime   time.Time
}

func (w Week) String() string {
	start := w.UnoffsetStart()
	end := w.UnoffsetEnd()
	return fmt.Sprintf("WK%d %d/%d-%d/%d", w.Index, start.Month(), start.Day(),
		end.Month(), end.Day())
}

func (w Week) StartDateStr() string {
	return FormatDate(w.StartTime.Add(-dayStart))
}

func (w Week) EndDateStr() string {
	return FormatDate(w.EndTime.Add(-dayStart - time.Nanosecond))
}

func (w Week) UnoffsetStart() time.Time {
	return w.StartTime.Add(-dayStart)
}

func (w Week) UnoffsetEnd() time.Time {
	return w.EndTime.Add(-dayStart - time.Nanosecond)
}

func (w Week) Equals(other Week) bool {
	return w.StartTime.Unix() == other.StartTime.Unix()
}

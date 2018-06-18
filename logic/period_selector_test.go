package logic

import (
	"testing"
	"time"
)

// TestInitPeriodSelector tests various aspects of the PeriodSelector
func TestInitPeriodSelector(t *testing.T) {
	periodStartDates := []time.Time{
		time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 1, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 2, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 3, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 4, 24, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 5, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 6, 19, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 7, 17, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 8, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 9, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 10, 9, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 11, 6, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 12, 4, 0, 0, 0, 0, time.UTC),
	}

	firstPeriod := []time.Time{
		time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	thirdPeriod := []time.Time{
		time.Date(2017, 1, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 2, 6, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 2, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2017, 2, 20, 0, 0, 0, 0, time.UTC),
	}

	periodSelector := NewPeriodSelector(2017)

	t.Run("A=1", func(t *testing.T) {
		for i, date := range periodStartDates {
			if periodSelector.Periods[i].StartTime() != date {
				t.Errorf("Dates do not match: %s != %s",
					periodSelector.Periods[i].StartTime().String(), date.String())
			}
		}
	})
	t.Run("A=2", func(t *testing.T) {
		for i, week := range periodSelector.Periods[0].Weeks {
			if firstPeriod[i] != week.StartTime {
				t.Errorf("Dates do not match: %s != %s",
					week.StartTime.String(), firstPeriod[i])
			}
		}

		for i, week := range periodSelector.Periods[2].Weeks {
			if thirdPeriod[i] != week.StartTime {
				t.Errorf("Dates do not match: %s != %s",
					week.StartTime.String(), thirdPeriod[i])
			}
		}
	})
}

func TestGetPeriodAndWeek(t *testing.T) {
	date := time.Date(2018, 3, 31, 0, 0, 0, 0, time.UTC)
	periodSelector := NewPeriodSelector(2018)

	t.Run("A=1", func(t *testing.T) {
		period := periodSelector.GetPeriod(date)

		if period.Index != 4 || period.StartTime().Day() != 26 {
			t.Errorf("Incorrect period: %s", period.String())
		}
	})
	t.Run("A=2", func(t *testing.T) {
		week := periodSelector.GetWeek(date)

		if week.Index != 1 || week.EndTime.Day() != 2 {
			t.Errorf("Incorrect week: %s", week.String())
		}
	})
	t.Run("A=3", func(t *testing.T) {
		period := periodSelector.GetPeriod(date)
		week := period.GetWeek(date)

		if week.Index != 1 || week.EndTime.Day() != 2 {
			t.Errorf("Incorrect week within period: %s", week.String())
		}
	})

}

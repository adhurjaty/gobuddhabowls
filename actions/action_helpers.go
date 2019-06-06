package actions

import (
	"buddhabowls/logic"
	"buddhabowls/presentation"
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

func setPeriodSelector(c buffalo.Context, presenter *presentation.Presenter) (time.Time, time.Time, error) {
	// get the parameters from URL
	paramsMap, ok := c.Params().(url.Values)
	if !ok {
		return time.Time{}, time.Time{}, c.Error(500, errors.New("Could not parse params"))
	}

	startVal, startTimeExists := paramsMap["StartTime"]
	endVal, endTimeExists := paramsMap["EndTime"]
	startTime := time.Time{}
	endTime := time.Time{}

	var err error
	if startTimeExists {
		startTime, err = time.Parse(time.RFC3339, startVal[0])
		if err != nil {
			return time.Time{}, time.Time{}, errors.WithStack(err)
		}
	} else if !endTimeExists {
		startTime = time.Now()
	}

	periods := presenter.GetPeriods(startTime)
	weeks := presenter.GetWeeks(startTime)
	years := presenter.GetYears()

	if endTimeExists {
		endTime, err = time.Parse(time.RFC3339, endVal[0])
		if err != nil {
			return time.Time{}, time.Time{}, errors.WithStack(err)
		}
		periods = append([]logic.Period{logic.Period{}}, periods...)
		weeks = append([]logic.Week{logic.Week{}}, weeks...)
		years = append([]int{0}, years...)

		c.Set("SelectedPeriod", periods[0])
		c.Set("SelectedWeek", weeks[0])
		c.Set("SelectedYear", startTime.Year())
	} else {
		selectedPeriod := presenter.GetSelectedPeriod(startTime)
		selectedWeek := presenter.GetSelectedWeek(startTime)
		c.Set("SelectedPeriod", selectedPeriod)
		c.Set("SelectedWeek", selectedWeek)
		c.Set("SelectedYear", startTime.Year())
		startTime = selectedWeek.StartTime
		endTime = selectedWeek.EndTime
	}

	c.Set("Periods", periods)
	c.Set("Weeks", weeks)
	c.Set("Years", years)
	c.Set("StartTime", startTime)
	c.Set("EndTime", endTime)

	return startTime, endTime, nil
}

package presentation

import (
	"buddhabowls/logic"
	"buddhabowls/models"
	"time"
)

// GetPeriods gets the list of periods available to the user
func (p *Presenter) GetPeriods(startTime time.Time) []logic.Period {
	if p.PeriodContext == nil {
		p.setPeriodContext(startTime)
	}

	return p.PeriodContext.PeriodSelector.Periods
}

// GetSelectedPeriod gets the period that contains the startTime
func (p *Presenter) GetSelectedPeriod(startTime time.Time) logic.Period {
	if p.PeriodContext == nil {
		p.setPeriodContext(startTime)
	}

	return p.PeriodContext.PeriodSelector.GetPeriod(startTime)
}

// GetWeeks gets the list of weeks available to the user
func (p *Presenter) GetWeeks(startTime time.Time) []logic.Week {
	if p.PeriodContext == nil {
		p.setPeriodContext(startTime)
	}

	return p.PeriodContext.SelectedPeriod.Weeks
}

// GetSelectedWeek gets the week that contains the startTime
func (p *Presenter) GetSelectedWeek(startTime time.Time) logic.Week {
	if p.PeriodContext == nil {
		p.setPeriodContext(startTime)
	}

	return p.PeriodContext.PeriodSelector.GetWeek(startTime)
}

// GetYears gets the list of years available to the user
func (p *Presenter) GetYears() []int {
	if p.PeriodContext == nil {
		years, err := models.GetYears(p.tx)
		if err != nil {
			return []int{1989}
		}
		return years
	}
	return p.PeriodContext.Years
}

// setPeriodContext sets a period selector context for using the period selector UI
func (p *Presenter) setPeriodContext(date time.Time) {
	p.PeriodContext = &periodSelectorContext{}
	p.PeriodContext.PeriodSelector = logic.NewPeriodSelector(date.Year())
	p.PeriodContext.SelectedPeriod = p.PeriodContext.PeriodSelector.GetPeriod(date)
	p.PeriodContext.SelectedWeek = p.PeriodContext.SelectedPeriod.GetWeek(date)
	p.PeriodContext.SelectedYear = date.Year()
	p.PeriodContext.Years, _ = models.GetYears(p.tx)
}

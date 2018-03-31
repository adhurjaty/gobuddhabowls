package componentcontexts

import (
	"buddhabowls/helpers"
	"time"
)

// PeriodSelectorContext context for the _period_selector.html template
type PeriodSelectorContext struct {
	PeriodSelector helpers.PeriodSelector
	SelectedYear   int
	SelectedPeriod helpers.Period
	SelectedWeek   helpers.Week
}

// Init initializes the context based on the date supplied
func (p *PeriodSelectorContext) Init(date time.Time) {
	p.PeriodSelector = *new(helpers.PeriodSelector)
	p.PeriodSelector.Init(date.Year())
	p.SelectedPeriod = p.PeriodSelector.GetPeriod(date)
	p.SelectedWeek = p.SelectedPeriod.GetWeek(date)
	p.SelectedYear = date.Year()
}

// func (p *PeriodSelectorContext) SetWeek(week helpers.Week) {
// 	p.SelectedWeek = week
// }

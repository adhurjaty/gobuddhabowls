package componentcontexts

import (
	"buddhabowls/logic"
	"time"
)

// PeriodSelectorContext context for the _period_selector.html template
type PeriodSelectorContext struct {
	PeriodSelector logic.PeriodSelector
	SelectedYear   int
	SelectedPeriod logic.Period
	SelectedWeek   logic.Week
}

// Init initializes the context based on the date supplied
func (p *PeriodSelectorContext) Init(date time.Time) {
	p.PeriodSelector = logic.NewPeriodSelector(date.Year())
	p.SelectedPeriod = p.PeriodSelector.GetPeriod(date)
	p.SelectedWeek = p.SelectedPeriod.GetWeek(date)
	p.SelectedYear = date.Year()
}

// func (p *PeriodSelectorContext) SetWeek(week helpers.Week) {
// 	p.SelectedWeek = week
// }

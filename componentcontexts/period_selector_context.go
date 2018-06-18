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

// NewPeriodSelectorContext initializes the context based on the date supplied
func NewPeriodSelectorContext(date time.Time) *PeriodSelectorContext {
	p := &PeriodSelectorContext{}
	p.PeriodSelector = logic.NewPeriodSelector(date.Year())
	p.SelectedPeriod = p.PeriodSelector.GetPeriod(date)
	p.SelectedWeek = p.SelectedPeriod.GetWeek(date)
	p.SelectedYear = date.Year()

	return p
}

// func (p *PeriodSelectorContext) SetWeek(week helpers.Week) {
// 	p.SelectedWeek = week
// }

package presentation

import (
	"buddhabowls/logic"
	"github.com/gobuffalo/pop"
)

// periodSelectorContext context for the _period_selector.html template
type periodSelectorContext struct {
	PeriodSelector logic.PeriodSelector
	SelectedYear   int
	SelectedPeriod logic.Period
	SelectedWeek   logic.Week
	Years          []int
}

// Presenter readies values for ViewModel consumption
type Presenter struct {
	tx            *pop.Connection
	PeriodContext *periodSelectorContext
}

func NewPresenter(tx *pop.Connection) *Presenter {
	p := &Presenter{tx, nil}
	return p
}

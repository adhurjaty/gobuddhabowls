package presentation

import (
	"buddhabowls/logic"
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"time"
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

// GetPeriodData gets all the period selector information to pass to the view
// func (p *Presenter) GetPeriodData(tx *pop.Connection) ([]logic.PeriodSelector, error) {
// 	years, err := models.GetYears(tx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	periods := make([]logic.PeriodSelector, len(years))

// 	for i, year := range years {
// 		periods[i] = logic.NewPeriodSelector(year)
// 	}

// 	return periods, nil
// }

// GetPurchaseOrders gets the purchase orders from the given date interval
func (p *Presenter) GetPurchaseOrders(startTime time.Time, endTime time.Time) (*PurchaseOrdersAPI, error) {
	purchaseOrders, err := logic.GetPurchaseOrders(startTime, endTime, p.tx)
	if err != nil {
		return nil, err
	}
	apiPurchaseOrders := &PurchaseOrdersAPI{}
	apiPurchaseOrders.ConvertToAPI(purchaseOrders)

	return apiPurchaseOrders, nil
}

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

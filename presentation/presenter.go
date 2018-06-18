package presentation

import (
	"buddhabowls/logic"
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"time"
)

// PeriodSelectorContext context for the _period_selector.html template
type PeriodSelectorContext struct {
	PeriodSelector logic.PeriodSelector
	SelectedYear   int
	SelectedPeriod logic.Period
	SelectedWeek   logic.Week
	Years          []int
}

// Presenter readies values for ViewModel consumption
type Presenter struct {
	tx *pop.Connection
}

func NewPresenter(tx *pop.Connection) *Presenter {
	p := &Presenter{tx}
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

// GetPeriodContext gets a period selector context for using the period selector UI
func (p *Presenter) GetPeriodContext(date time.Time) *PeriodSelectorContext {
	ps := &PeriodSelectorContext{}
	ps.PeriodSelector = logic.NewPeriodSelector(date.Year())
	ps.SelectedPeriod = ps.PeriodSelector.GetPeriod(date)
	ps.SelectedWeek = ps.SelectedPeriod.GetWeek(date)
	ps.SelectedYear = date.Year()
	ps.Years, _ = models.GetYears(p.tx)

	return ps
}

// GetPurchaseOrders gets the purchase orders from the given date interval
func (p *Presenter) GetPurchaseOrders(startTime time.Time, endTime time.Time) (*models.PurchaseOrders, error) {
	q := p.tx.Eager().Where("order_date >= ? AND order_date < ?",
		startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)).Order("order_date DESC")

	factory := models.ModelFactory{}
	pos := &models.PurchaseOrders{}
	err := factory.CreateModelSlice(pos, q)

	return pos, err
}

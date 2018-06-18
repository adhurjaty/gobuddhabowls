package presentation

import (
	"buddhabowls/componentcontexts"
	"buddhabowls/logic"
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"time"
)

// Presenter readies values for ViewModel consumption
type Presenter struct{}

// GetPeriodData gets all the period selector information to pass to the view
func (p *Presenter) GetPeriodData(tx *pop.Connection) ([]logic.PeriodSelector, error) {
	years, err := models.GetYears(tx)
	if err != nil {
		return nil, err
	}
	periods := make([]logic.PeriodSelector, len(years))

	for i, year := range years {
		periods[i] = logic.NewPeriodSelector(year)
	}

	return periods, nil
}

// GetPeriodContext gets a period selector context for using the period selector UI
func (p *Presenter) GetPeriodContext(date time.Time) *componentcontexts.PeriodSelectorContext {
	return componentcontexts.NewPeriodSelectorContext(date)
}

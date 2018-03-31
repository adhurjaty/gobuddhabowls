package componentcontexts

import (
	"testing"
	"time"
)

func TestInitPeriodSelectorContext(t *testing.T) {
	date := time.Date(2018, 3, 31, 0, 0, 0, 0, time.UTC)
	pSelectorContext := PeriodSelectorContext{}
	pSelectorContext.Init(date)

	t.Run("A=1", func(t *testing.T) {
		if pSelectorContext.SelectedPeriod.StartDateStr() != "03/26/2018" {
			t.Errorf("Incorrect selected period: %s", pSelectorContext.SelectedPeriod.String())
		}
		if pSelectorContext.SelectedWeek.StartDateStr() != "03/26/2018" {
			t.Errorf("Incorrect selected week: %s", pSelectorContext.SelectedWeek.String())
		}
	})
	t.Run("A=2", func(t *testing.T) {
		refWeeks := []string{
			"03/26/2018",
			"04/02/2018",
			"04/09/2018",
			"04/16/2018",
		}
		if len(pSelectorContext.SelectedPeriod.Weeks) != 4 {
			t.Errorf("Incorrect number of weeks in period")
		}
		for i, week := range pSelectorContext.SelectedPeriod.Weeks {
			if refWeeks[i] != week.StartDateStr() {
				t.Errorf("Incorrect week within period: %s", week.String())
			}
		}
	})
}

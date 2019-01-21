package actions

import (
	"buddhabowls/presentation"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func ListSales(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	_, _, err := setPeriodSelector(c, presenter)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.HTML("sales/index"))
}

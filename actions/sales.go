package actions

import (
	"github.com/gobuffalo/buffalo"
)

func ListSales(c buffalo.Context) error {

	return c.Render(200, r.HTML("sales/index"))
}

package actions

import (
	"github.com/gobuffalo/buffalo"
)

// SettingsHandler serves the main settings page
// GET /settings
func SettingsHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("settings/index"))
}

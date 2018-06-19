package presentation

// CategoryAPI category object for ui
type CategoryAPI struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Background string `json:"background"`
	Index      int    `json:"index"`
}

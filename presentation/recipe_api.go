package presentation

type RecipeAPI struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	IsMenu   string      `json:"is_menu"`
	Category CategoryAPI `json:"Category"`
	Items    ItemsAPI    `json:"Items"`
	Index    int         `json:"index"`
}

type RecipesAPI []RecipeAPI

func (r RecipeAPI) String() string {
	jo, _ := json.Marshal(r)
	return string(jo)
}

func (r RecipesAPI) String() string {
	jo, _ := json.Marshal(r)
	return string(jo)
}

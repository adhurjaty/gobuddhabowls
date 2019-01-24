package presentation

import (
	"buddhabowls/models"
	"encoding/json"
)

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

func NewRecipesAPI(recipes *models.Recipes) RecipesAPI {
	apis := RecipesAPI{}
	for _, recipe := range *recipes {
		apis = append(apis, NewRecipeAPI(&recipe))
	}

	return apis
}

func NewRecipeAPI(recipe *models.Recipe) RecipeAPI {
	return RecipeAPI{
		ID:       recipe.ID.String(),
		Name:     recipe.Name,
		Category: NewCategoryAPI(recipe.Category),
		Index:    recipe.Index,
		Items:    NewItemsAPI(recipe.Items),
	}
}

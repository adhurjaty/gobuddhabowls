package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"github.com/gobuffalo/uuid"
)

type RecipeAPI struct {
	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	IsBatch              bool        `json:"is_batch"`
	Category             CategoryAPI `json:"Category"`
	RecipeUnit           string      `json:"recipe_unit"`
	RecipeUnitConversion float64     `json:"recipe_unit_conversion"`
	Items                ItemsAPI    `json:"Items"`
	Index                int         `json:"index"`
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
		ID:                   recipe.ID.String(),
		Name:                 recipe.Name,
		Category:             NewCategoryAPI(recipe.Category),
		Index:                recipe.Index,
		RecipeUnit:           recipe.RecipeUnit,
		RecipeUnitConversion: recipe.RecipeUnitConversion,
		Items:                NewItemsAPI(recipe.Items),
		IsBatch:              recipe.IsBatch,
	}
}

func ConvertToModelRecipe(recAPI *RecipeAPI) (*models.Recipe, error) {
	id, err := uuid.FromString(recAPI.ID)
	if err != nil {
		id = uuid.UUID{}
	}

	// filter out 0 count items
	recItems := ItemsAPI{}
	for _, item := range recAPI.Items {
		if item.Count > 0 {
			recItems = append(recItems, item)
		}
	}
	items, err := ConvertToModelRecipeItems(recItems, id)
	if err != nil {
		return nil, err
	}
	category, err := ConvertToModelCategory(recAPI.Category)
	if err != nil {
		return nil, err
	}

	return &models.Recipe{
		ID:                   id,
		Name:                 recAPI.Name,
		CategoryID:           category.ID,
		Category:             *category,
		RecipeUnit:           recAPI.RecipeUnit,
		RecipeUnitConversion: recAPI.RecipeUnitConversion,
		IsBatch:              recAPI.IsBatch,
		Items:                *items,
		Index:                recAPI.Index,
	}, nil
}

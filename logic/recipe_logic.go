package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
)

func GetRecipes(tx *pop.Connection) (*models.Recipes, error) {
	factory := models.ModelFactory{}
	recipes := &models.Recipes{}
	if err := factory.CreateModelSlice(recipes, tx.Eager().Q()); err != nil {
		return nil, err
	}

	recipes.Sort()

	return recipes, nil
}

func GetRecipe(id string, tx *pop.Connection) (*models.Recipe, error) {
	factory := models.ModelFactory{}
	recipe := &models.Recipe{}
	err := factory.CreateModel(recipe, tx, id)

	return recipe, err
}

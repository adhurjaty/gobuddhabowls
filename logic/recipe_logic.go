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

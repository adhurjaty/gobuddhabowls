package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
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

func GetBatchRecipes(tx *pop.Connection) (*models.Recipes, error) {
	factory := models.ModelFactory{}
	recipes := &models.Recipes{}
	query := tx.Eager().Where("is_batch = true")
	if err := factory.CreateModelSlice(recipes, query); err != nil {
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

func UpdateRecipe(recipe *models.Recipe, tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	oldRecipe, err := GetRecipe(recipe.ID.String(), tx)
	if err != nil {
		return verrs, err
	}
	oldItems := oldRecipe.Items

	containsFunc := func(item models.RecipeItem, itemArr models.RecipeItems) bool {
		for _, otherItem := range itemArr {
			if item.ID == otherItem.ID {
				return true
			}
		}
		return false
	}

	for _, item := range recipe.Items {
		item.RecipeID = recipe.ID
		if containsFunc(item, oldItems) {
			verrs, err = tx.ValidateAndUpdate(&item)
		} else {
			verrs, err = tx.ValidateAndCreate(&item)
		}
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	// delete items removed from recipe
	for _, item := range oldItems {
		if !containsFunc(item, recipe.Items) {
			err = tx.Destroy(&item)
			if err != nil {
				return verrs, err
			}
		}
	}

	return UpdateRecipeNoItems(recipe, tx)
}

func UpdateRecipeNoItems(recipe *models.Recipe, tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndUpdate(recipe)
}

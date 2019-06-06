package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

func GetRecipes(tx *pop.Connection) (*models.Recipes, error) {
	return getRecipesHelper(tx.Eager().Q())
}

func GetRecipesNoItems(tx *pop.Connection) (*models.Recipes, error) {
	recipes := &models.Recipes{}
	err := tx.Eager().All(recipes)
	recipes.Sort()
	return recipes, err
}

func GetBatchRecipes(tx *pop.Connection) (*models.Recipes, error) {
	query := tx.Eager().Where("is_batch = true")
	return getRecipesHelper(query)
}

func GetRecipesOfCategory(id string, catID string, tx *pop.Connection) (*models.Recipes, error) {
	query := tx.Eager().Where("category_id = ?", catID).Where("id != ?", id)
	return getRecipesHelper(query)
}

func getRecipesHelper(query *pop.Query) (*models.Recipes, error) {
	factory := models.ModelFactory{}
	recipes := &models.Recipes{}
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
	verrs, err := updateRecIndices(recipe, tx)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}
	return tx.ValidateAndUpdate(recipe)
}

func updateRecIndices(recItem *models.Recipe, tx *pop.Connection) (*validate.Errors, error) {
	items, err := GetRecipesOfCategory(recItem.ID.String(),
		recItem.CategoryID.String(), tx)
	if err != nil {
		return validate.NewErrors(), err
	}

	offset := 0
	for i, item := range *items {
		if item.Index == recItem.Index {
			offset = 1
		}

		item.Index = i + offset
		verrs, err := tx.ValidateAndUpdate(&item)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return validate.NewErrors(), nil
}

func InsertRecipe(recipe *models.Recipe, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := tx.ValidateAndCreate(recipe)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	for _, item := range recipe.Items {
		item.RecipeID = recipe.ID
		verrs, err = tx.ValidateAndCreate(&item)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return validate.NewErrors(), nil
}

func DestroyRecipe(recipe *models.Recipe, tx *pop.Connection) error {
	for _, item := range recipe.Items {
		if err := tx.Destroy(&item); err != nil {
			return err
		}
	}

	return tx.Destroy(recipe)
}

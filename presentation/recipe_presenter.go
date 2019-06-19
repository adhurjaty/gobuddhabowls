package presentation

import (
	"buddhabowls/logic"

	"github.com/gobuffalo/validate"
)

func (p *Presenter) GetRecipe(id string) (*RecipeAPI, error) {
	recipe, err := logic.GetRecipe(id, p.tx)
	if err != nil {
		return nil, err
	}

	recAPI := NewRecipeAPI(recipe)
	err = p.populateRecipeItemCosts(&recAPI.Items)
	return &recAPI, err
}

func (p *Presenter) GetRecipes() (*RecipesAPI, error) {
	recipes, err := logic.GetRecipes(p.tx)
	if err != nil {
		return nil, err
	}

	recipesAPI := NewRecipesAPI(recipes)

	for i, _ := range recipesAPI {
		rec := &recipesAPI[i]
		err = p.populateRecipeItemCosts(&rec.Items)
		if err != nil {
			return nil, err
		}
	}

	return &recipesAPI, nil
}

func (p *Presenter) GetBatchRecipes() (*RecipesAPI, error) {
	recipes, err := logic.GetBatchRecipes(p.tx)
	if err != nil {
		return nil, err
	}

	recipesAPI := NewRecipesAPI(recipes)
	return &recipesAPI, nil
}

func (p *Presenter) populateRecipeItemCosts(items *ItemsAPI) error {
	for i := 0; i < len(*items); i++ {
		cost, err := p.getItemRecipeCost((*items)[i])
		if err != nil {
			return err
		}
		(*items)[i].Price = cost
	}

	return nil
}

func (p *Presenter) getItemRecipeCost(item ItemAPI) (float64, error) {
	cost := 0.0
	if item.InventoryItemID != "" {
		vendorItem, err := logic.GetSelectedVendorItem(item.InventoryItemID, p.tx)
		if err != nil {
			// if there's no matching vendor item - just say price is 0
			return 0, nil
		}
		return vendorItem.Price / vendorItem.Conversion /
			item.RecipeUnitConversion, nil
	}

	recipe, err := logic.GetRecipe(item.BatchRecipeID, p.tx)
	if err != nil {
		return 0, err
	}
	recAPI := NewRecipeAPI(recipe)
	for _, rItem := range recAPI.Items {
		incCost, err := p.getItemRecipeCost(rItem)
		if err != nil {
			return 0, err
		}
		cost += incCost * rItem.Count
	}

	return cost / recAPI.RecipeUnitConversion, nil
}

func (p *Presenter) GetRecipesNoItems() (*RecipesAPI, error) {
	recipes, err := logic.GetRecipesNoItems(p.tx)
	if err != nil {
		return nil, err
	}

	recipesAPI := NewRecipesAPI(recipes)
	return &recipesAPI, nil
}

func (p *Presenter) GetAllItemsForRecipe() (*ItemsAPI, error) {
	batchRecipes, err := logic.GetBatchRecipes(p.tx)
	if err != nil {
		return nil, err
	}

	batchItems := NewItemsAPI(batchRecipes)

	items, err := p.GetInventoryItems()
	if err != nil {
		return nil, err
	}

	*items = append(*items, batchItems...)
	err = p.populateRecipeItemCosts(items)
	clearItemIds(items)

	return items, err
}

func (p *Presenter) UpdateRecipe(recAPI *RecipeAPI) (*validate.Errors, error) {
	recipe, err := ConvertToModelRecipe(recAPI)
	if err != nil {
		return validate.NewErrors(), err
	}

	return logic.UpdateRecipe(recipe, p.tx)
}

func (p *Presenter) UpdateRecipeNoItems(recAPI *RecipeAPI) (*validate.Errors, error) {
	recipe, err := ConvertToModelRecipe(recAPI)
	if err != nil {
		return validate.NewErrors(), err
	}

	return logic.UpdateRecipeNoItems(recipe, p.tx)
}

func (p *Presenter) InsertRecipe(recAPI *RecipeAPI) (*validate.Errors, error) {
	recipe, err := ConvertToModelRecipe(recAPI)
	if err != nil {
		return validate.NewErrors(), err
	}

	return logic.InsertRecipe(recipe, p.tx)
}

func (p *Presenter) DestroyRecipe(recAPI *RecipeAPI) error {
	recipe, err := ConvertToModelRecipe(recAPI)
	if err != nil {
		return err
	}

	return logic.DestroyRecipe(recipe, p.tx)
}

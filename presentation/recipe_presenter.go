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
	err = p.populateReciepItemCosts(&recAPI.Items)
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
		err = p.populateReciepItemCosts(&rec.Items)
		if err != nil {
			return nil, err
		}
	}

	return &recipesAPI, nil
}

func (p *Presenter) populateReciepItemCosts(items *ItemsAPI) error {
	for i := range *items {
		item := &(*items)[i]
		cost, err := p.getItemRecipeCost(item)
		if err != nil {
			return err
		}
		item.Price = cost
	}

	return nil
}

func (p *Presenter) getItemRecipeCost(item *ItemAPI) (float64, error) {
	cost := 0.0
	if item.InventoryItemID != "" {
		vendorItem, err := logic.GetSelectedVendorItem(item.InventoryItemID, p.tx)
		if err != nil {
			// if there's no matching vendor item - just say price is 0
			return 0, nil
		}
		cost = vendorItem.Price / vendorItem.Conversion /
			item.RecipeUnitConversion
	} else {
		recipe, err := logic.GetRecipe(item.BatchRecipeID, p.tx)
		if err != nil {
			return 0, err
		}
		recAPI := NewRecipeAPI(recipe)
		for _, item := range recAPI.Items {
			incCost, err := p.getItemRecipeCost(&item)
			if err != nil {
				return 0, err
			}
			cost += incCost * item.Count
		}

		cost /= recAPI.RecipeUnitConversion
	}

	return cost, nil
}

func (p *Presenter) GetAllItemsForRecipe() (*ItemsAPI, error) {
	batchRecipes, err := logic.GetBatchRecipes(p.tx)
	if err != nil {
		return nil, err
	}

	batchItems := NewItemsAPI(*batchRecipes)

	items, err := p.GetInventoryItems()
	if err != nil {
		return nil, err
	}

	*items = append(*items, batchItems...)
	err = p.populateReciepItemCosts(items)

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

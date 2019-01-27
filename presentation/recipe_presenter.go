package presentation

import (
	"buddhabowls/logic"
)

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
	for i, _ := range *items {
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
	vendorItem, err := logic.GetSelectedVendorItem(item.InventoryItemID, p.tx)
	cost := 0.0
	if err == nil {
		cost = vendorItem.Price / vendorItem.Conversion /
			item.RecipeUnitConversion * item.Count
	} else {
		recipe, err := logic.GetRecipe(item.InventoryItemID, p.tx)
		if err != nil {
			return 0, err
		}
		recAPI := NewRecipeAPI(recipe)
		for _, item := range recAPI.Items {
			incCost, err := p.getItemRecipeCost(&item)
			if err != nil {
				return 0, err
			}
			cost += incCost
		}
	}

	return cost, nil
}

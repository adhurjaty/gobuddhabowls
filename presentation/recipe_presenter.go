package presentation

func (p *Presenter) GetRecipes() (*RecipesAPI, error) {
	recipes, err := logic.GetRecipes(p.tx)
	if err != nil {
		return nil, err
	}

	recipeAPI := NewRecipesAPI(recipes)

	return &recipeAPI, nil
}

package presentation

import (
	"buddhabowls/logic"
)

func (p *Presenter) GetAllCategories() (*CategoriesAPI, error) {
	categories, err := logic.GetAllCategories(p.tx)
	if err != nil {
		return nil, err
	}

	catAPI := NewCategoriesAPI(categories)
	return &catAPI, nil
}

func (p *Presenter) GetInvItemCategories() (*CategoriesAPI, error) {
	categories, err := logic.GetInvItemCategories(p.tx)
	if err != nil {
		return nil, err
	}

	catAPI := NewCategoriesAPI(categories)
	return &catAPI, nil
}

func (p *Presenter) GetRecCategories() (*CategoriesAPI, error) {
	categories, err := logic.GetRecCategories(p.tx)
	if err != nil {
		return nil, err
	}

	catAPI := NewCategoriesAPI(categories)
	return &catAPI, nil
}

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
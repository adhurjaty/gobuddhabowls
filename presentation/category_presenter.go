package presentation

import (
	"buddhabowls/logic"
	"github.com/gobuffalo/validate"
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

func (p *Presenter) InsertCategory(catAPI *CategoryAPI) (*validate.Errors, error) {
	category, err := ConvertToModelCategory(*catAPI)
	if err != nil {
		return validate.NewErrors(), err
	}

	verrs, err := logic.InsertCategory(category, p.tx)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}
	catAPI.ID = category.ID.String()
	return validate.NewErrors(), nil
}

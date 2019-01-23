package presentation

import (
	"buddhabowls/logic"
	"github.com/gobuffalo/validate"
)

func (p *Presenter) GetUser(id string) (*UserAPI, error) {
	user, err := logic.GetUser(id, p.tx)
	if err != nil {
		return nil, err
	}

	userAPI := NewUserAPI(user)
	return &userAPI, nil
}

func (p *Presenter) UpdateUser(userAPI *UserAPI) (*validate.Errors, error) {
	user, err := ConvertToModelUser(userAPI)
	if err != nil {
		return nil, err
	}

	return logic.UpdateUser(user, p.tx)
}

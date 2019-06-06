package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"github.com/gobuffalo/uuid"
)

type UserAPI struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	PasswordHash   string `json:"-"`
	SquareLocation string `json:"square_location"`
	SquareToken    string `json:"square_token"`
	Timezone       string `json:"timezone"`
}

func (u UserAPI) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

func NewUserAPI(user *models.User) UserAPI {
	return UserAPI{
		ID:             user.ID.String(),
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		SquareLocation: user.SquareLocation.String,
		SquareToken:    user.SquareToken.String,
		Timezone:       user.Timezone,
	}
}

func ConvertToModelUser(userAPI *UserAPI) (*models.User, error) {
	id, err := uuid.FromString(userAPI.ID)
	if err != nil {
		id = uuid.UUID{}
	}

	return &models.User{
		ID:             id,
		Email:          userAPI.Email,
		PasswordHash:   userAPI.PasswordHash,
		SquareLocation: models.StringToNullString(userAPI.SquareLocation),
		SquareToken:    models.StringToNullString(userAPI.SquareToken),
		Timezone:       userAPI.Timezone,
	}, nil
}

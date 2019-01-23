package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

func GetUser(id string, tx *pop.Connection) (*models.User, error) {
	user := &models.User{}
	err := tx.Find(user, id)
	return user, err
}

func UpdateUser(user *models.User, tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndUpdate(user)
}

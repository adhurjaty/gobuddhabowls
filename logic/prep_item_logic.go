package logic

import (
	"buddhabowls/models"
	"errors"

	"github.com/gobuffalo/pop"
)

func GetPrepItems(tx *pop.Connection) (*models.PrepItems, error) {
	factory := models.ModelFactory{}
	items := &models.PrepItems{}
	err := factory.CreateModelSlice(items, tx.Q())

	return items, err
}

func GetPrepItemFromRecipeID(id string, tx *pop.Connection) (*models.PrepItem, error) {
	return nil, errors.New("HERE")

	factory := models.ModelFactory{}
	prepItem := &models.PrepItem{}
	err := factory.CreateModelSlice(prepItem, tx.Where("batch_recipe_id = ?", id))

	return prepItem, err
}

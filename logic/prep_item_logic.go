package logic

import (
	"buddhabowls/models"

	"github.com/gobuffalo/pop"
)

func GetPrepItems(tx *pop.Connection) (*models.PrepItems, error) {
	factory := models.ModelFactory{}
	items := &models.PrepItems{}
	err := factory.CreateModelSlice(items, tx.Q())

	return items, err
}
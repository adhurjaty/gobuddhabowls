package logic

import (
	"buddhabowls/models"
	"fmt"

	"github.com/gobuffalo/pop"
)

func GetPrepItems(tx *pop.Connection) (*models.PrepItems, error) {
	factory := models.ModelFactory{}
	items := &models.PrepItems{}
	err := factory.CreateModelSlice(items, tx.Q())

	return items, err
}

func GetPrepItemFromRecipeID(id string, tx *pop.Connection) (*models.PrepItem, error) {
	factory := models.ModelFactory{}
	prepItems := &models.PrepItems{}
	if err := factory.CreateModelSlice(prepItems, tx.Where("batch_recipe_id = ?", id)); err != nil {
		return nil, err
	}
	if len(*prepItems) == 0 {
		return nil, fmt.Errorf("no prep item with recipe ID: %s", id)
	}
	fmt.Println(prepItems)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	return &(*prepItems)[0], nil
}

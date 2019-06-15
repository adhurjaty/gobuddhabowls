package logic

import (
	"buddhabowls/models"
	"fmt"

	"github.com/gobuffalo/validate"

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
	return &(*prepItems)[0], nil
}

func GetPrepItem(id string, tx *pop.Connection) (*models.PrepItem, error) {
	factory := models.ModelFactory{}
	prepItem := &models.PrepItem{}
	err := factory.CreateModel(prepItem, tx, id)

	return prepItem, err
}

func UpdatePrepItem(prepItem *models.PrepItem, tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := UpdateIndices(prepItem, tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return tx.ValidateAndUpdate(prepItem)
}

func GetPrepItemsOfCategory(prepItem *models.PrepItem, tx *pop.Connection) (*models.PrepItems, error) {
	query := tx.RawQuery(fmt.Sprintf(
		`SELECT pi.* FROM prep_items AS pi
		JOIN recipes AS r ON r.id = pi.batch_recipe_id
		WHERE r.category_id IN 
			(SELECT r.category_id FROM prep_items AS pi
			JOIN recipes AS r ON pi.batch_recipe_id = r.id
			WHERE pi.id = '%s')`, prepItem.ID.String()))
	return getPrepItemsHelper(query)
}

func getPrepItemsHelper(query *pop.Query) (*models.PrepItems, error) {
	factory := models.ModelFactory{}
	items := &models.PrepItems{}
	err := factory.CreateModelSlice(items, query)

	if err != nil {
		return nil, err
	}

	items.Sort()

	return items, err
}

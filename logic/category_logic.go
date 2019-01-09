package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"sort"
)

func GetAllCategories(tx *pop.Connection) (*models.InventoryItemCategories, error) {
	categories := &models.InventoryItemCategories{}
	if err := tx.Eager().All(categories); err != nil {
		return nil, err
	}

	sort.Slice(*categories, func(i, j int) bool {
		return (*categories)[i].Index < (*categories)[j].Index
	})

	return categories, nil
}

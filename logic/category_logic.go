package logic

import (
	"buddhabowls/models"
	"github.com/gobuffalo/pop"
	"sort"
)

func GetAllCategories(tx *pop.Connection) (*models.ItemCategories, error) {
	return getCategoriesFromQuery(tx.Q())
}

func GetInvItemCategories(tx *pop.Connection) (*models.ItemCategories, error) {
	query := tx.Eager().Where()
	return getCategoriesFromQuery(query)
}

func GetRecCategories(tx *pop.Connection) (*models.ItemCategories, error) {
	query := tx.Eager().Where()
	return getCategoriesFromQuery(query)
}

func getCategoriesFromQuery(query *pop.Query) (*models.ItemCategories, error) {
	categories := &models.ItemCategories{}
	if err := query.All(categories); err != nil {
		return nil, err
	}

	sort.Slice(*categories, func(i, j int) bool {
		return (*categories)[i].Index < (*categories)[j].Index
	})

	return categories, nil
}

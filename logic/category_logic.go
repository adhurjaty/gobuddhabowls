package logic

import (
	"buddhabowls/models"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"sort"
)

func GetAllCategories(tx *pop.Connection) (*models.ItemCategories, error) {
	return getCategoriesFromQuery(tx.Q())
}

func GetInvItemCategories(tx *pop.Connection) (*models.ItemCategories, error) {
	queryString := createCategoryTypeQueryString("inventory_items")
	query := tx.Eager().RawQuery(queryString)
	return getCategoriesFromQuery(query)
}

func GetRecCategories(tx *pop.Connection) (*models.ItemCategories, error) {
	queryString := createCategoryTypeQueryString("recipes")
	query := tx.Eager().RawQuery(queryString)
	return getCategoriesFromQuery(query)
}

func createCategoryTypeQueryString(tableName string) string {
	return fmt.Sprintf("SELECT DISTINCT ci.* FROM item_categories AS ci"+
		" INNER JOIN %s AS t ON ci.id = t.category_id", tableName)
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

func InsertCategory(category *models.ItemCategory, tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(category)
}

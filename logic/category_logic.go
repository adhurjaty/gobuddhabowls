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

func GetAllRecCategories(tx *pop.Connection) (*models.RecipeCategories, error) {
	categories := &models.RecipeCategories{}
	if err := tx.Eager().All(categories); err != nil {
		return nil, err
	}

	sort.Slice(*categories, func(i, j int) bool {
		return (*categories)[i].Index < (*categories)[j].Index
	})

	return categories, nil
}

func InvCategoryIntSlice(categories *models.InventoryItemCategories) *models.Categories {
	outCats := &models.Categories{}
	for _, cat := range *categories {
		*outCats = append(*outCats, cat)
	}

	return outCats
}

func RecCategoryIntSlice(categories *models.RecipeCategories) *models.Categories {
	outCats := &models.Categories{}
	for _, cat := range *categories {
		*outCats = append(*outCats, cat)
	}

	return outCats
}

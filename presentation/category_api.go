package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"github.com/gobuffalo/uuid"
)

// CategoryAPI category object for ui
type CategoryAPI struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Background string `json:"background"`
	Index      int    `json:"index"`
}

type CategoriesAPI []CategoryAPI

func (c CategoryAPI) String() string {
	jo, _ := json.Marshal(c)
	return string(jo)
}

// NewCategoryAPI converts a category to an api category
func NewCategoryAPI(category models.Category) CategoryAPI {
	c := CategoryAPI{}

	c.ID = category.GetID().String()
	c.Name = category.GetName()
	c.Background = category.GetBackground()
	c.Index = category.GetIndex()

	return c
}

func NewCategoriesAPI(categories *models.Categories) CategoriesAPI {
	catsAPI := CategoriesAPI{}
	for _, category := range *categories {
		catsAPI = append(catsAPI, NewCategoryAPI(category))
	}

	return catsAPI
}

func ConvertToModelCategory(catAPI CategoryAPI) (*models.ItemCategory, error) {
	id, err := uuid.FromString(catAPI.ID)
	if err != nil {
		return nil, err
	}

	return &models.ItemCategory{
		ID:         id,
		Name:       catAPI.Name,
		Background: catAPI.Background,
		Index:      catAPI.Index,
	}, nil
}

// SelectValue returns the ID for select input tags
func (c CategoryAPI) SelectValue() interface{} {
	return c.ID
}

// SelectLabel returs the name for select input tags
func (c CategoryAPI) SelectLabel() string {
	if c.ID == "" {
		return ""
	}
	return c.Name
}

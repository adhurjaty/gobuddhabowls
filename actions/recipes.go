package actions

import (
	"buddhabowls/models"
	"buddhabowls/presentation"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

var _ = fmt.Print

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Recipe)
// DB Table: Plural (recipes)
// Resource: Plural (Recipes)
// Path: Plural (/recipes)
// View Template Folder: Plural (/templates/recipes/)

// RecipesResource is the resource for the Recipe model
type RecipesResource struct {
	buffalo.Resource
}

// List gets all Recipes. This function is mapped to the path
// GET /recipes
func (v RecipesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	recipes, err := presenter.GetRecipes()
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println(recipes)

	c.Set("recipes", recipes)

	return c.Render(200, r.HTML("recipes/index"))
}

// Show gets the data for one Recipe. This function is mapped to
// the path GET /recipes/{recipe_id}
func (v RecipesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Recipe
	recipe := &models.Recipe{}

	// To find the Recipe the parameter recipe_id is used.
	if err := tx.Find(recipe, c.Param("recipe_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, recipe))
}

// New renders the form for creating a new Recipe.
// This function is mapped to the path GET /recipes/new
func (v RecipesResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Recipe{}))
}

// Create adds a Recipe to the DB. This function is mapped to the
// path POST /recipes
func (v RecipesResource) Create(c buffalo.Context) error {
	// Allocate an empty Recipe
	recipe := &models.Recipe{}

	// Bind recipe to the html form elements
	if err := c.Bind(recipe); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(recipe)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, recipe))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Recipe was created successfully")

	// and redirect to the recipes index page
	return c.Render(201, r.Auto(c, recipe))
}

// Edit renders a edit form for a Recipe. This function is
// mapped to the path GET /recipes/{recipe_id}/edit
func (v RecipesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Recipe
	recipe := &models.Recipe{}

	if err := tx.Find(recipe, c.Param("recipe_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, recipe))
}

// Update changes a Recipe in the DB. This function is mapped to
// the path PUT /recipes/{recipe_id}
func (v RecipesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Recipe
	recipe := &models.Recipe{}

	if err := tx.Find(recipe, c.Param("recipe_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Recipe to the html form elements
	if err := c.Bind(recipe); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(recipe)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, recipe))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Recipe was updated successfully")

	// and redirect to the recipes index page
	return c.Render(200, r.Auto(c, recipe))
}

// Destroy deletes a Recipe from the DB. This function is mapped
// to the path DELETE /recipes/{recipe_id}
func (v RecipesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Recipe
	recipe := &models.Recipe{}

	// To find the Recipe the parameter recipe_id is used.
	if err := tx.Find(recipe, c.Param("recipe_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(recipe); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Recipe was destroyed successfully")

	// Redirect to the recipes index page
	return c.Render(200, r.Auto(c, recipe))
}

package actions

import (
	"buddhabowls/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func UsersNew(c buffalo.Context) error {
	u := models.User{}
	c.Set("user", u)
	return c.Render(200, r.HTML("users/new.html"))
}

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(200, r.HTML("users/new.html"))
	}

	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "Welcome to Buffalo!")

	return c.Redirect(302, "/")
}

// GET /users/square
func UsersSquare(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	user := &models.User{}
	if err := tx.Find(user, c.Session().Get("current_user_id")); err != nil {
		return errors.WithStack(err)
	}

	c.Set("user", user)

	return c.Render(200, r.HTML("users/square_settings"))
}

// PUT to /users/square
func UpdateUsersSquare(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	user := &models.User{}
	if err := tx.Find(user, c.Session().Get("current_user_id")); err != nil {
		return errors.WithStack(err)
	}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)

		return c.Render(200, r.HTML("users/square_settings"))
	}

	c.Flash().Add("success", "Successfully updated Square settings")

	return c.Redirect(303, "/sales")
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/signin")
		}
		return next(c)
	}
}

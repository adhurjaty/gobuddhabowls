package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/mw-forcessl"
	"github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	"buddhabowls/models"

	"github.com/gobuffalo/mw-csrf"
	"github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_buddhabowls_session",
		})
		// Automatically redirect to SSL
		app.Use(forcessl.Middleware(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))
		app.Use(SetCurrentUser)
		app.Use(Authorize)

		app.Middleware.Skip(Authorize, AuthCreate, AuthNew)

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)
		app.GET("/routes", RoutesHandler)

		app.PUT("/inventory_items/{inventory_item_id}/inline", UpdateInventoryItem)
		app.Resource("/inventory_items", InventoryItemsResource{})

		app.Resource("/inventories", InventoriesResource{})

		app.PUT("/vendors/{vendor_id}/inline", UpdateVendorInline)
		app.Resource("/vendors", VendorsResource{})
		app.Resource("/count_inventory_items", CountInventoryItemsResource{})

		app.GET("/purchase_orders/{purchase_order_id}/receiving_list", ShowReceivingList)
		app.GET("/purchase_orders/{purchase_order_id}/order_sheet", ShowOrderSheet)
		app.Resource("/purchase_orders", PurchaseOrdersResource{})

		app.GET("/sales", ListSales)

		app.GET("/users/square", UsersSquare)
		app.PUT("/users/square", UpdateUsersSquare)

		app.Resource("/order_items", OrderItemsResource{})
		app.Resource("/prep_items", PrepItemsResource{})

		app.PUT("/recipes/{recipe_id}/inline", UpdateRecipeInline)
		app.Resource("/recipes", RecipesResource{})

		app.Resource("/recipe_items", RecipeItemsResource{})
		app.Resource("/vendor_items", VendorItemsResource{})
		app.Resource("/count_prep_items", CountPrepItemsResource{})

		app.GET("/settings", SettingsHandler)
		app.Resource("/inventory_item_categories", InventoryItemCategoriesResource{})
		app.Resource("/recipe_categories", RecipeCategoriesResource{})

		app.GET("/signin", AuthNew)
		app.POST("/signin", AuthCreate)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

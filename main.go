// main.go
package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/pocketbase/pocketbase/core"
	"github.com/supercip971/ulist-backend/pkg/urequests"

	_ "github.com/supercip971/ulist-backend/migrations"
)

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true, // auto creates migration files when making collection changes
	})
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/version",
			Handler: func(c echo.Context) error {
				return c.String(http.StatusOK, "1")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		// api/v1/list-entries
		urequests.AddListEntriesQueriesRoutes(e, app)

		// (post) api/v1/list-entries
		urequests.AddListEntriesUpdatesRoutes(e, app)

		// api/v1/list-join/:id
		urequests.AddInvitationJoiningRoutes(e, app)

		// (post) api/v1/list-invite
		urequests.AddInvitationCreationRoutes(e, app)

		// (post) api/v1/list
		urequests.AddListCreationRoute(e, app)

		// api/v1/list
		urequests.AddListInfoRoutes(e, app)

		// api/v1/list-properties
		urequests.AddListPropertiesRoutes(e, app)

		// api/v1/get-lists
		urequests.AddListQueriesRoutes(e, app)
		// Getting user lists API

		return nil
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

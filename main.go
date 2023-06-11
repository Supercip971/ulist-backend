/*
 Copyright (C) 2023 cyp

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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

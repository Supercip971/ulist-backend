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

package urequests

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/supercip971/ulist-backend/pkg/db"
)

func AddListCreationRoute(e *core.ServeEvent, app *pocketbase.PocketBase) {

	// List Create API
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/list",
		Handler: func(c echo.Context) error {
			l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

			if l == nil {
				return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
			}

			params := db.CreateShoppingList{}

			err1 := c.Bind(&params)

			if err1 != nil {
				return c.JSON(http.StatusBadRequest, "{\"error\":\"unable to create list\"}")
			}
			err := DbListCreate(app, params, l.GetId())

			if err != nil {

				return c.JSON(http.StatusForbidden, "{\"error\":\"there was an issue creating this list\"}")
			}

			return c.JSON(http.StatusOK, "{}")

		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})
}

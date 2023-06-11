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
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/supercip971/ulist-backend/pkg/db"
)

func AddListPropertiesRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/list-properties",
		Handler: func(c echo.Context) error {
			l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

			if l == nil {
				return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
			}

			userId := l.GetId()

			query := app.Dao().ModelQuery(&db.ListVisibility{})

			result := []*db.ListVisibility{}
			err := query.AndWhere(dbx.HashExp{"userOwnership": userId}).All(&result)

			if err != nil {

				return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
			}

			returned := db.UserLists{
				Count: len(result),
			}

			for _, v := range result {
				ent := db.UserListEntry{
					Id:      v.Id,
					ListId:  v.List,
					UserId:  v.UserOwnership,
					Owner:   v.Owner,
					Created: v.Created.String(),
				}
				returned.Lists = append(returned.Lists, ent)
			}

			return c.JSON(http.StatusOK, returned)

		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})
}

func AddListInfoRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/list",
		Handler: func(c echo.Context) error {
			l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

			if l == nil {
				return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
			}
			if !c.QueryParams().Has("id") {
				return c.JSON(http.StatusForbidden, "{\"error\":\"Please enter a valid 'id'\"}")
			}
			id := c.QueryParams().Get("id")

			userId := l.GetId()
			if !DoUserHasRights(app, userId, id) {
				return c.JSON(http.StatusForbidden, "{\"error\":\"The user doesn't have access to this list.\"}")
			}

			result := &db.ShopList{}
			err := app.Dao().ModelQuery(&db.ShopList{}).AndWhere(dbx.HashExp{"id": id}).One(result)
			//	err := query.AndWhere(dbx.HashExp{"userOwnership": userId}).All(&result)

			if err != nil {

				return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
			}

			return c.JSON(http.StatusOK, result)

		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})
}

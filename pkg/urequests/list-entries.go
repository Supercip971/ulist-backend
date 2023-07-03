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
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/supercip971/ulist-backend/pkg/db"
)

func AddListEntriesQueriesRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	// List entries API
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/list-entries",
		Handler: func(c echo.Context) error {
			// sleep
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

			result := []*db.DbShopListEntries{}
			err := app.Dao().ModelQuery(&db.DbShopListEntries{}).AndWhere(dbx.HashExp{"list": id}).All(&result)
			//	err := query.AndWhere(dbx.HashExp{"userOwnership": userId}).All(&result)

			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusForbidden, "{\"error\":\"An error has occured during the query\"}")
			}

			returned := db.ShoppingList{
				Count: len(result),
			}

			for _, v := range result {
				username, err := app.Dao().FindRecordById("users", v.AddedBy)
				if err != nil {
					print("unable to found user: ", v.AddedBy)
					print(err)
					continue
				}
				
				ent := db.ShoppingListEntry{
					Name:           v.Name,
					AddedBy:        username.Username(),
					Checked:        v.Checked,
					CustomRelation: v.CustomRelation,
					Id:             v.Id,
					List:           v.List,
					Tags: 			v.Tags,
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

func AddListEntriesUpdatesRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	// List entries API (update)
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/list-entries",
		Handler: func(c echo.Context) error {
			l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

			if l == nil {
				return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
			}

			params := db.PostShoppingList{}

			err1 := c.Bind(&params)

			//val := []byte{}

			//				c.Request().Body.Read(val)

			//				print("query string: '" + c.QueryString() + "'\n")

			//				print("query body: '" + string(val) + "'\n")
			//		print("query body: '" + (c.Get("data").(string)) + "'\n")

			//				err1 := json.Unmarshal(val, &params)

			if err1 != nil {
				//		print("invalid query: " + err1.Error() + "req: " + string(val))
				return c.JSON(http.StatusForbidden, "{\"error\":\"An error has occured during the query\"}")
			}

			v, _ := json.MarshalIndent(params, "", "  ")
			println("params: " + string(v))
			if !c.QueryParams().Has("id") {
				return c.JSON(http.StatusForbidden, "{\"error\":\"Please enter a valid 'id'\"}")
			}
			id := c.QueryParams().Get("id")

			// maybe useless, but we never know
			if id != params.ListId {
				return c.JSON(http.StatusForbidden, "{\"error\":\"Please enter a valid 'id'\"}")
			}

			userId := l.GetId()
			if !DoUserHasRights(app, userId, id) {
				return c.JSON(http.StatusForbidden, "{\"error\":\"The user doesn't have access to this list.\"}")
			}

			err := DbListUpdate(app, params, userId)

			//	err := query.AndWhere(dbx.HashExp{"userOwnership": userId}).All(&result)

			if err != nil {
				print("error during db update: " + err.Error())
				return c.JSON(http.StatusForbidden, "{\"error\":\"An error has occured during the query\"}")
			}

			return c.JSON(http.StatusOK, "")

		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})
}

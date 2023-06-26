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
				return c.JSON(
					http.StatusForbidden,
					"{\"error\":\"forbidden, only connected user has the right to use this api\"}",
				)
			}

			userId := l.GetId()

			if !c.QueryParams().Has("id") {
				return c.JSON(http.StatusForbidden, "{\"error\":\"Please enter a valid 'id'\"}")
			}
			id := c.QueryParams().Get("id")

			if !DoUserHasRights(app, userId, id) {
				return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, user doesn't have the right to use this list \"}")
			}
			returned := db.ShoppingListProperties{
				Users:  []db.GetUserInList{},
				Shares: []db.GetListShare{},
			}

			/* users in list */

			query := app.Dao().ModelQuery(&db.ListVisibility{})
			users := []*db.ListVisibility{}
			err := query.AndWhere(dbx.HashExp{"list": id}).All(&users)
			if err != nil {
				print("unable to get users in the list")
				print(err)
				return c.JSON(http.StatusInternalServerError, "{\"error\":\"Unable to get users in the list\"}")
			}
			for _, user := range users {

				username, err := app.Dao().FindRecordById("users", user.UserOwnership)
				if err != nil {
					print("unable to found user: ", user.UserOwnership)
					print(err)
					continue
				}

				returned_user := db.GetUserInList{
					Name:  username.Username(),
					Id:    user.Id,
					Owner: user.Owner,
				}
				returned.Users = append(returned.Users, returned_user)
			}

			/* shares in list */

			share_query := app.Dao().ModelQuery(&db.ListShare{})
			shares := []*db.ListShare{}
			err = share_query.AndWhere(dbx.HashExp{"list": id}).All(&shares)

			if err != nil { 

				print("unable to get shares in the list", err)
				return c.JSON(http.StatusInternalServerError, "{\"error\":\"Unable to get shares in the list\"}")
			}

			for _, share := range shares {
				returned_share := db.GetListShare{
					Identifier:     share.Identifier,
					SharedBy:       share.SharedBy,
					ExpirationDate: share.ExpirationDate,
				}

				returned.Shares = append(returned.Shares, returned_share)
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
				return c.JSON(
					http.StatusForbidden,
					"{\"error\":\"forbidden, only connected user has the right to use this api\"}",
				)
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
				return c.JSON(
					http.StatusForbidden,
					"{\"error\":\"forbidden, only connected user has the right to use this api\"}",
				)
			}

			return c.JSON(http.StatusOK, result)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})
}

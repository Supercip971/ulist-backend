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

func AddListQueriesRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/get-lists",
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

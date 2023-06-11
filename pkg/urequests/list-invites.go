package urequests

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tools/types"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/supercip971/ulist-backend/pkg/db"
)

func AddInvitationJoiningRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/list-join/:id",
		Handler: func(c echo.Context) error {

			invite_id := c.PathParam("id")

			if invite_id == "" {

				return echo.NewHTTPError(http.StatusBadRequest, "Please enter a valid 'id'")
				//	return c.JSON(http.StatusBadRequest, "{\"error\":\"Please enter a valid 'id'\"}")
			}

			l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

			if l == nil {
				return echo.NewHTTPError(http.StatusForbidden, "Forbidden, only connected user has the right to use this api")

			}

			userId := l.GetId()

			if !UserAntispamCheck(userId) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "You are trying to join too many lists at the same time, please wait 5-10 seconds.")

			}

			//	if !do_user_has_right(app, userId, id) {
			//		return c.JSON(http.StatusForbidden, "{\"error\":\"The user doesn't have access to this list.\"}")
			//	}

			result := &db.ListShare{}
			err := app.Dao().ModelQuery(&db.ListShare{}).AndWhere(dbx.HashExp{"id": invite_id}).One(result)

			//	err := query.AndWhere(dbx.HashExp{"userOwnership": userId}).All(&result)

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid invite id, maybe the invitation has expired or the 'id' is invalid.")

			}

			if !result.ExpirationDate.IsZero() {
				if result.ExpirationDate.Time().Sub(time.Now()) < 0 {
					return echo.NewHTTPError(http.StatusBadRequest, "The invitation has expired.")

				}
			}

			println("user joining list id: " + userId)
			visibility := &db.ListVisibility{
				List:          result.List,
				UserOwnership: userId,
				Owner:         false,
				Readonly:      false,
			}

			visibility.MarkAsNew()
			err = app.Dao().Save(visibility)

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, result)

		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})
}

func AddInvitationCreationRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	// List Invitation Create API
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/list-invite",
		Handler: func(c echo.Context) error {
			l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

			if l == nil {
				return echo.NewHTTPError(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
			}

			params := db.ListShare{}

			err1 := c.Bind(&params)

			if err1 != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "{\"error\":\"unable to create invitation\"}")
			}

			if params.ExpirationDate.IsZero() || params.ExpirationDate.Time().Sub(time.Now().Add(time.Hour*24*7)) < 0 {

				params.ExpirationDate, _ = types.ParseDateTime(time.Now().Add(time.Hour * 24 * 7).UTC().Format(types.DefaultDateLayout))
			}
			params.SharedBy = l.GetId()
			params.MarkAsNew()
			err := app.Dao().Save(&params)

			params.Identifier = params.Id
			if err != nil {

				return echo.NewHTTPError(http.StatusForbidden, "{\"error\":\"there was an issue creating this invitation\"}")
			}

			return c.JSON(http.StatusOK, params.Id)

		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})
}

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

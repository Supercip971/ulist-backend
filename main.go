// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/types"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/supercip971/ulist-backend/pkg/db"

	_ "github.com/supercip971/ulist-backend/migrations"
)

// ensures that the List struct satisfy the models.Model interface

type AntiSpamEntry struct {
	time time.Time
	end  time.Time
}

var antiSpam = make(map[string]AntiSpamEntry)

func userAntispamUpdate(userId string) bool {

	for key, value := range antiSpam {

		// now > end
		// now - end > 0
		if time.Now().Sub(value.end) > time.Second*5 {
			delete(antiSpam, key)
		}
	}

	return false
}

func userAntispamCheck(userId string) bool {

	entry, ok := antiSpam[userId]

	if ok {

		// now > start + 10s
		if time.Now().Sub(entry.time) > time.Second*10 {
			delete(antiSpam, userId)
			return true
		}

	} else {
		antiSpam[userId] = AntiSpamEntry{
			time: time.Now(),
			end:  time.Now().Add(time.Second * 10),
		}
		return true
	}
	return false
}

func do_user_has_right(app *pocketbase.PocketBase, userId string, listId string) bool {
	query := app.Dao().ModelQuery(&db.ListVisibility{})

	result := []*db.ListVisibility{}
	err := query.AndWhere(dbx.HashExp{"userOwnership": userId, "list": listId}).All(&result)

	return err == nil && len(result) > 0
}

func db_list_create(app *pocketbase.PocketBase, req db.CreateShoppingList, userId string) error {

	record := &db.ShopList{
		Name:              req.Name,
		ReadonlyByDefault: false,
	}

	record.MarkAsNew()
	err := app.Dao().Save(record)

	if err != nil {
		return err
	}

	println("record id: " + record.Id)
	visibility := &db.ListVisibility{
		List:          record.Id,
		UserOwnership: userId,
		Owner:         true,
		Readonly:      false,
	}

	visibility.MarkAsNew()
	err = app.Dao().Save(visibility)

	if err != nil {
		return err
	}

	return nil
}

func db_list_update(app *pocketbase.PocketBase, req db.PostShoppingList, userId string) error {

	//collection, err := app.Dao().FindCollectionByNameOrId(&db.DbShopListEntries{}.TableName())

	// do all creation before the update
	for _, entry := range req.Entries {

		if entry.Name == "" {
			continue
		}
		res := &db.DbShopListEntries{
			Name:    entry.Name,
			List:    req.ListId,
			AddedBy: userId,
			Checked: entry.Checked,
		}
		res.MarkAsNew()
		if err := app.Dao().Save(res); err != nil {
			return err
		}

		//		err = app.Dao().SaveRecord(record)
	}
	for _, entry := range req.Update {

		original := &db.DbShopListEntries{}

		err := app.Dao().ModelQuery(&db.DbShopListEntries{}).AndWhere(dbx.HashExp{"id": entry.EntryId}).One(original)

		if err != nil {
			return err
		}

		if entry.Name != "" {
			original.Name = entry.Name
		}
		original.Checked = entry.Checked

		original.MarkAsNotNew()
		if err := app.Dao().Save(original); err != nil {
			return err
		}

		//		err = app.Dao().SaveRecord(record)
	}

	return nil
}

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true, // auto creates migration files when making collection changes
	})
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/add-member",
			Handler: func(c echo.Context) error {

				return c.String(http.StatusOK, "Hello world!")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth(),
			},
		})

		// List entries API
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/v1/list-entries",
			Handler: func(c echo.Context) error {
				// sleep
				time.Sleep(300 * time.Millisecond)
				l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

				if l == nil {
					return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
				}
				if !c.QueryParams().Has("id") {
					return c.JSON(http.StatusForbidden, "{\"error\":\"Please enter a valid 'id'\"}")
				}
				id := c.QueryParams().Get("id")

				userId := l.GetId()
				if !do_user_has_right(app, userId, id) {
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
					ent := db.ShoppingListEntry{
						Name:           v.Name,
						AddedBy:        v.AddedBy,
						Checked:        v.Checked,
						CustomRelation: v.CustomRelation,
						Id:             v.Id,
						List:           v.List,
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
				if !do_user_has_right(app, userId, id) {
					return c.JSON(http.StatusForbidden, "{\"error\":\"The user doesn't have access to this list.\"}")
				}

				err := db_list_update(app, params, userId)

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

		// List Invitation Join API
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

					return c.JSON(http.StatusBadRequest, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
				}

				userId := l.GetId()

				if userAntispamCheck(userId) {
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
		// List Invitation Create API
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/v1/list-invite",
			Handler: func(c echo.Context) error {
				l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

				if l == nil {
					return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
				}

				params := db.ListShare{}

				err1 := c.Bind(&params)

				if err1 != nil {
					return c.JSON(http.StatusBadRequest, "{\"error\":\"unable to create invitation\"}")
				}

				if params.ExpirationDate.IsZero() || params.ExpirationDate.Time().Sub(time.Now().Add(time.Hour*24*7)) < 0 {

					params.ExpirationDate, _ = types.ParseDateTime(time.Now().Add(time.Hour * 24 * 7).UTC().Format(types.DefaultDateLayout))
				}
				params.MarkAsNew()
				err := app.Dao().Save(&params)

				if err != nil {

					return c.JSON(http.StatusForbidden, "{\"error\":\"there was an issue creating this invitation\"}")
				}

				return c.JSON(http.StatusOK, params.Id)

			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth(),
			},
		})
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
				err := db_list_create(app, params, l.GetId())

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
		// List information API
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
				if !do_user_has_right(app, userId, id) {
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

		// Getting user lists API
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
		return nil
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

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
	"os"

	"time"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/supercip971/ulist-backend/pkg/db"

	"github.com/stripe/stripe-go/sub"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"database/sql"
)


const (
	SubNull = -1
	SubFree = 0
	SubPremium = 1
)


var StripeEnabled = false
var StripeSubscriptionKey = ""
func EnsureDbStripeUser(app *pocketbase.PocketBase,  user *models.Record) (res error, val db.UserSubscriptionInformation) {
	
	previous := db.UserSubscriptionInformation{}

	query := app.Dao().ModelQuery(&db.UserSubscriptionInformation{})

	err := query.AndWhere(dbx.HashExp{"user": user.GetId()}).One(&previous)
	

	if err != nil {
		res = nil 
		val = previous 
		return 
	}

	if err != sql.ErrNoRows { 
		res = err
		return 
	}	
	

	params := &stripe.CustomerParams{
		Description: stripe.String("Customer for ulist"),
		Email: stripe.String(user.Email()), 
		PreferredLocales: []*string{stripe.String("en")},
	} 

	c, err1 := customer.New(params)

	if(err1 != nil) {
		res = err1
	}

	print("loaded user stripe information: ", c.ID)

	sub_info := db.UserSubscriptionInformation{
		UserId: user.GetId(),
		StripeId: c.ID,
		SubscriptionTier: 0,
	}
	sub_info.MarkAsNew()
	err2 := app.Dao().Save(&sub_info)


	if(err2 != nil) {
		res = err2
		return 
	}

	res = nil 
	val = sub_info

	return 
}


func UserPremiumLevel(app *pocketbase.PocketBase, userId string) int {
	
	// for selfhosting, not using a STRIPE key will make the app think that all users are premium
	if(!StripeEnabled) {
		return SubPremium
	}
	sub_info := db.UserSubscriptionInformation{}

	query := app.Dao().ModelQuery(&db.UserSubscriptionInformation{})

	err := query.AndWhere(dbx.HashExp{"user": userId}).One(&sub_info)

	
	if err != nil {
		println("no payment information for user"); 
		return SubFree
	}


	if  time.Now().Sub(sub_info.Updated.Time()).Minutes() < 1 {

		return sub_info.SubscriptionTier
	}
	st_id := sub_info.StripeId
	
	println("stripe id: ", st_id)

	params := &stripe.SubscriptionListParams{
		Customer: st_id,
	}

	params.Filters.AddFilter("limit", "", "1")
	
	i := sub.List(params)
	if(i == nil) {

		print("no subscription for user")
		return SubFree 

	}


	for i.Next() {

		s := i.Subscription()

		if s.Plan.Product.ID == StripeSubscriptionKey {

			sub_info.SubscriptionTier = SubPremium


			app.Dao().Save(&sub_info)

			return SubPremium
		}

		println("ff: {}", s.Plan.Product.ID)
		
	}
	sub_info.SubscriptionTier = SubFree

	app.Dao().Save(&sub_info)


	print("no matched subscription for user")

	return SubFree 
}

func StripeInit(){
	key := os.Getenv("STRIPE_SECRET_KEY")
	if(key == ""){

		println("STRIPE_SECRET_KEY not set, stripe will not be available")
		return

	}
	stripe.Key = key
	StripeEnabled = true
	StripeSubscriptionKey = os.Getenv("STRIPE_SUBSCRIPTION_KEY")
	println("Stripe is enabled")
}
func AddListStripeRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/subscription",
		Handler: func(c echo.Context) error {


			l, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

			if l == nil {
				return c.JSON(http.StatusForbidden, "{\"error\":\"forbidden, only connected user has the right to use this api\"}")
			}

			userId := l.GetId()
			
			level := UserPremiumLevel(app, userId)

			return c.JSON(http.StatusOK, level)


		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})
}

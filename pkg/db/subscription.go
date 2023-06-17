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

package db

import (
	"github.com/pocketbase/pocketbase/models"
)

type UserSubscriptionInformation struct {
	models.BaseModel

	UserId 		  	 string   `db:"user" json:"user"`
	StripeId         string   `db:"stripeId" json:"stripeId"`
	SubscriptionTier int      `db:"subTier" json:"subTier"`
}

var _ models.Model = (*UserSubscriptionInformation)(nil)

func (m *UserSubscriptionInformation) TableName() string {
	return "subscription_info" // the name of your collection
}


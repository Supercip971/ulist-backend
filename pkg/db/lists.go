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
	"github.com/pocketbase/pocketbase/tools/types"
)

type ListVisibility struct {
	models.BaseModel

	UserOwnership string `db:"userOwnership" json:"userOwnership"`
	Owner         bool   `db:"owner" json:"owner"`
	Readonly      bool   `db:"readonly" json:"readonly"`
	List          string `db:"list" json:"list"`
}

var _ models.Model = (*ListVisibility)(nil)

func (m *ListVisibility) TableName() string {
	return "lists_visibility" // the name of your collection
}

type ShopList struct {
	models.BaseModel

	Name              string `db:"name" json:"name"`
	ReadonlyByDefault bool   `db:"readonly_by_default" json:"readonly_by_default"`
}

var _ models.Model = (*ShopList)(nil)

func (m *ShopList) TableName() string {
	return "lists" // the name of your collection
}

type DbShopListEntries struct {
	models.BaseModel

	Name string `db:"name" json:"name"`

	List           string `db:"list" json:"list"`
	AddedBy        string `db:"addedBy" json:"addedBy"`
	Checked        bool   `db:"checked" json:"checked"`
	CustomRelation string `db:"customRelation" json:"customRelation"`
}

var _ models.Model = (*DbShopListEntries)(nil)

func (m *DbShopListEntries) TableName() string {
	return "list_entries" // the name of your collection
}

type ListShare struct {
	models.BaseModel

	List           string         `db:"list" json:"list"`
	SharedBy       string         `db:"sharedBy" json:"sharedBy"`
	ExpirationDate types.DateTime `db:"expirationDate" json:"expirationDate"`
	Identifier     string         `db:"identificator" json:"identificator"`
}

var _ models.Model = (*ListShare)(nil)

func (m *ListShare) TableName() string {
	return "list_shares" // the name of your collection
}

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

import "github.com/pocketbase/pocketbase/tools/types"


type ShoppingListEntry struct {
	Id             string `db:"id" json:"id"`
	Created        string `db:"created" json:"created"`
	Name           string `db:"name" json:"name"`
	List           string `db:"list" json:"list"`
	AddedBy        string `db:"addedBy" json:"addedBy"`
	Checked        bool   `db:"checked" json:"checked"`
	Tags 		   string `db:"tags" json:"tags"`
	CustomRelation string `db:"customRelation" json:"customRelation"`
}

type ShoppingList struct {
	Lists []ShoppingListEntry `json:"entries"`
	Count int                 `json:"count"`
}


type UserListEntry struct {
	Id      string `db:"id" json:"id"`
	Created string `db:"created" json:"created"`
	ListId  string `db:"listId" json:"listId"`
	UserId  string `db:"userId" json:"userId"`
	Owner   bool   `db:"owner" json:"owner"`
}

type CreateShoppingList struct {
	Name string `json:"name"`
}

type UserLists struct {
	Lists []UserListEntry `json:"lists"`
	Count int             `json:"count"`
}

type PutShoppingList struct {
	Name    string `db:"name" json:"name"`
	Checked bool   `db:"checked" json:"checked"`
	Tags    string `db:"tags" json:"tags"`
}

type UpdateShoppingList struct {
	EntryId string `db:"id" json:"id"`     // entry id, using 0 for new entry
	Name    string `db:"name" json:"name"` //
	Checked bool   `db:"checked" json:"checked"`
	Tags 	string `db:"tags" json:"tags"`
	Delete  bool   `db:"delete" json:"delete"`
}

type PostShoppingList struct {
	ListId string `db:"id" json:"id"`

	Entries []PutShoppingList    `db:"pushes" json:"pushes"`
	Update  []UpdateShoppingList `db:"updates" json:"updates"`
}

type GetListShare struct {
	SharedBy       string         `db:"sharedBy" json:"sharedBy"`
	ExpirationDate types.DateTime `db:"expirationDate" json:"expirationDate"`
	Identifier     string         `db:"identificator" json:"identificator"`
}

type GetUserInList struct {
	Name string `json:"name"`
	Id string `json:"id"`
	Owner bool `json:"owner"`
}

type ShoppingListProperties struct {
	Name 		string `json:"name"`

	Users 		[]GetUserInList `json:"users"`
	Shares 		[]GetListShare `json:"shares"`
}

}

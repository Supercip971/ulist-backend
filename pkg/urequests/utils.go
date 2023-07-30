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
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/supercip971/ulist-backend/pkg/db"
)

// ensures that the List struct satisfy the models.Model interface

type AntiSpamEntry struct {
	time time.Time
	end  time.Time
}

var antiSpam = make(map[string]AntiSpamEntry)

func UserAntispamUpdate(userId string) bool {

	for key, value := range antiSpam {

		// now > end
		// now - end > 0
		if time.Now().Sub(value.end) > time.Second*5 {
			delete(antiSpam, key)
		}
	}

	return false
}

func UserAntispamCheck(userId string) bool {

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

func DoUserHasRights(app *pocketbase.PocketBase, userId string, listId string) bool {
	query := app.Dao().ModelQuery(&db.ListVisibility{})

	result := []*db.ListVisibility{}
	err := query.AndWhere(dbx.HashExp{"userOwnership": userId, "list": listId}).All(&result)

	return err == nil && len(result) > 0
}

func DbListCreate(app *pocketbase.PocketBase, req db.CreateShoppingList, userId string) error {

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

func DbListUpdate(app *pocketbase.PocketBase, req db.PostShoppingList, userId string) error {

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
			Tags:   entry.Tags,
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


		if(entry.Delete) {
			if err := app.Dao().Delete(original); err != nil {
				return err
			}
			continue
		}

		if entry.Name != "" {
			original.Name = entry.Name
		}
		original.Checked = entry.Checked

		if entry.Tags != "NONE" {
			original.Tags = entry.Tags
		}

		original.MarkAsNotNew()
		if err := app.Dao().Save(original); err != nil {
			return err
		}

		//		err = app.Dao().SaveRecord(record)
	}

	return nil
}

package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ew68jtyrjuejmap")
		if err != nil {
			return err
		}

		// update
		edit_expirationDate := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "d4kdwssp",
			"name": "expirationDate",
			"type": "date",
			"required": true,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), edit_expirationDate)
		collection.Schema.AddField(edit_expirationDate)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ew68jtyrjuejmap")
		if err != nil {
			return err
		}

		// update
		edit_expirationDate := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "d4kdwssp",
			"name": "expiration_date",
			"type": "date",
			"required": true,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), edit_expirationDate)
		collection.Schema.AddField(edit_expirationDate)

		return dao.SaveCollection(collection)
	})
}

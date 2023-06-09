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
		edit_expiration_date := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "d4kdwssp",
			"name": "expiration_date",
			"type": "date",
			"required": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), edit_expiration_date)
		collection.Schema.AddField(edit_expiration_date)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ew68jtyrjuejmap")
		if err != nil {
			return err
		}

		// update
		edit_expiration_date := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "d4kdwssp",
			"name": "expirationDate",
			"type": "date",
			"required": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), edit_expiration_date)
		collection.Schema.AddField(edit_expiration_date)

		return dao.SaveCollection(collection)
	})
}

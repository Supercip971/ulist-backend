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

		// add
		new_identificator := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "hc9yflf1",
			"name": "identificator",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_identificator)
		collection.Schema.AddField(new_identificator)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ew68jtyrjuejmap")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("hc9yflf1")

		return dao.SaveCollection(collection)
	})
}

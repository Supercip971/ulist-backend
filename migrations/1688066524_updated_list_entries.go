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

		collection, err := dao.FindCollectionByNameOrId("yynw9pxmd71jghn")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("jgi5nv8e")

		// add
		new_tags := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "orwqotjt",
			"name": "tags",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_tags)
		collection.Schema.AddField(new_tags)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("yynw9pxmd71jghn")
		if err != nil {
			return err
		}

		// add
		del_tags := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "jgi5nv8e",
			"name": "tags",
			"type": "json",
			"required": false,
			"unique": false,
			"options": {}
		}`), del_tags)
		collection.Schema.AddField(del_tags)

		// remove
		collection.Schema.RemoveField("orwqotjt")

		return dao.SaveCollection(collection)
	})
}

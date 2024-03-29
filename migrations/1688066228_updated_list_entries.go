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

		// add
		new_tags := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "jgi5nv8e",
			"name": "tags",
			"type": "json",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_tags)
		collection.Schema.AddField(new_tags)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("yynw9pxmd71jghn")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("jgi5nv8e")

		return dao.SaveCollection(collection)
	})
}

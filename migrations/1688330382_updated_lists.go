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

		collection, err := dao.FindCollectionByNameOrId("42bpfxq80b2axz8")
		if err != nil {
			return err
		}

		// add
		new_tags := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "yhay62ms",
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

		collection, err := dao.FindCollectionByNameOrId("42bpfxq80b2axz8")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("yhay62ms")

		return dao.SaveCollection(collection)
	})
}

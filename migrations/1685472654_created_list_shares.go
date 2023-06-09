package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "ew68jtyrjuejmap",
			"created": "2023-05-30 18:50:54.715Z",
			"updated": "2023-05-30 18:50:54.715Z",
			"name": "list_shares",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "ng4f4bun",
					"name": "list",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "42bpfxq80b2axz8",
						"cascadeDelete": true,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "8mbickuw",
					"name": "sharedBy",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": true,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
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
				}
			],
			"indexes": [],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ew68jtyrjuejmap")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}

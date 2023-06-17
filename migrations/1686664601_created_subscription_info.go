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
			"id": "dff7ai1t2ji827r",
			"created": "2023-06-13 13:56:41.662Z",
			"updated": "2023-06-13 13:56:41.662Z",
			"name": "subscription_info",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "pmj027od",
					"name": "user",
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
					"id": "qugdum5o",
					"name": "stripeId",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "uylw7uhm",
					"name": "subTier",
					"type": "number",
					"required": false,
					"unique": false,
					"options": {
						"min": 0,
						"max": null
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

		collection, err := dao.FindCollectionByNameOrId("dff7ai1t2ji827r")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}

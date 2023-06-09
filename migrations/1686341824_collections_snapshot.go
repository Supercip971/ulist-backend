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
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2023-05-01 15:18:13.731Z",
				"updated": "2023-05-30 20:06:48.913Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
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
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"protected": false
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"requireEmail": false
				}
			},
			{
				"id": "42bpfxq80b2axz8",
				"created": "2023-05-07 20:08:35.274Z",
				"updated": "2023-05-30 20:06:48.917Z",
				"name": "lists",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "budcgyi9",
						"name": "name",
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
						"id": "k4kr5sgc",
						"name": "readonly_by_default",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "8wa8r2k2epw10ci",
				"created": "2023-05-07 20:10:03.584Z",
				"updated": "2023-05-30 20:06:48.917Z",
				"name": "lists_visibility",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "vbzpjfti",
						"name": "userOwnership",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "eg0swzaz",
						"name": "owner",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "p7lxu5qb",
						"name": "readonly",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "6fwuhfeh",
						"name": "list",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "42bpfxq80b2axz8",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
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
			},
			{
				"id": "yynw9pxmd71jghn",
				"created": "2023-05-07 20:11:23.333Z",
				"updated": "2023-05-30 20:06:48.917Z",
				"name": "list_entries",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "t52cb2lo",
						"name": "name",
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
						"id": "nguhpogl",
						"name": "list",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "42bpfxq80b2axz8",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "ivmnxdnb",
						"name": "addedBy",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "l4acwtvu",
						"name": "checked",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "yzatsjl1",
						"name": "customRelation",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
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
			},
			{
				"id": "fakw5n71js0v6of",
				"created": "2023-05-09 20:06:13.372Z",
				"updated": "2023-05-30 20:06:48.917Z",
				"name": "emoji_prompt",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "41zc57o0",
						"name": "prompt",
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
						"id": "48h6fuh4",
						"name": "emoji",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
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
			},
			{
				"id": "ew68jtyrjuejmap",
				"created": "2023-05-30 18:50:54.715Z",
				"updated": "2023-06-09 19:23:23.831Z",
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
						"name": "expirationDate",
						"type": "date",
						"required": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
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
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}

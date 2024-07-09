package pb_migrations

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

		collection, err := dao.FindCollectionByNameOrId("s1q7t7ofpbuozf9")
		if err != nil {
			return err
		}

		// update
		edit_tg_username := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "v0rqqqev",
			"name": "tg_username",
			"type": "text",
			"required": true,
			"presentable": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_tg_username); err != nil {
			return err
		}
		collection.Schema.AddField(edit_tg_username)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("s1q7t7ofpbuozf9")
		if err != nil {
			return err
		}

		// update
		edit_tg_username := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "v0rqqqev",
			"name": "tg_username",
			"type": "text",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_tg_username); err != nil {
			return err
		}
		collection.Schema.AddField(edit_tg_username)

		return dao.SaveCollection(collection)
	})
}

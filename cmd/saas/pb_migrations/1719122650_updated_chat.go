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

		// add
		new_linked_chat_id := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "2talvlsr",
			"name": "linked_chat_id",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "s1q7t7ofpbuozf9",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), new_linked_chat_id); err != nil {
			return err
		}
		collection.Schema.AddField(new_linked_chat_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("s1q7t7ofpbuozf9")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("2talvlsr")

		return dao.SaveCollection(collection)
	})
}

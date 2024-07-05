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

		collection, err := dao.FindCollectionByNameOrId("52sylu6udk1kc6r")
		if err != nil {
			return err
		}

		// add
		new_is_tg_history_message := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "xes0wuad",
			"name": "is_tg_history_message",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_is_tg_history_message); err != nil {
			return err
		}
		collection.Schema.AddField(new_is_tg_history_message)

		// add
		new_tg_history_entities := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "4vydhvuf",
			"name": "tg_history_entities",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), new_tg_history_entities); err != nil {
			return err
		}
		collection.Schema.AddField(new_tg_history_entities)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("52sylu6udk1kc6r")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("xes0wuad")

		// remove
		collection.Schema.RemoveField("4vydhvuf")

		return dao.SaveCollection(collection)
	})
}

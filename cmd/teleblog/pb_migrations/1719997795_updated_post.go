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

		// remove
		collection.Schema.RemoveField("pn9i4rge")

		// remove
		collection.Schema.RemoveField("4vydhvuf")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("52sylu6udk1kc6r")
		if err != nil {
			return err
		}

		// add
		del_tg_entities := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "pn9i4rge",
			"name": "tg_entities",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), del_tg_entities); err != nil {
			return err
		}
		collection.Schema.AddField(del_tg_entities)

		// add
		del_tg_history_entities := &schema.SchemaField{}
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
		}`), del_tg_history_entities); err != nil {
			return err
		}
		collection.Schema.AddField(del_tg_history_entities)

		return dao.SaveCollection(collection)
	})
}

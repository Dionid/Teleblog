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

		collection, err := dao.FindCollectionByNameOrId("f7ecawbcx0paa90")
		if err != nil {
			return err
		}

		// add
		new_is_tg_history_message := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "cfvycymp",
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

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("f7ecawbcx0paa90")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("cfvycymp")

		return dao.SaveCollection(collection)
	})
}

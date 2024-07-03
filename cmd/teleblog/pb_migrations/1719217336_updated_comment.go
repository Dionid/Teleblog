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

		// remove
		collection.Schema.RemoveField("slxi4mjm")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("f7ecawbcx0paa90")
		if err != nil {
			return err
		}

		// add
		del_reply_to_id := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "slxi4mjm",
			"name": "reply_to_id",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "f7ecawbcx0paa90",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), del_reply_to_id); err != nil {
			return err
		}
		collection.Schema.AddField(del_reply_to_id)

		return dao.SaveCollection(collection)
	})
}

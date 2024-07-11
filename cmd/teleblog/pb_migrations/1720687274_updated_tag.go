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

		collection, err := dao.FindCollectionByNameOrId("2bepntx0gwpms2d")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("kdkvxytx")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("2bepntx0gwpms2d")
		if err != nil {
			return err
		}

		// add
		del_post_id := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "kdkvxytx",
			"name": "post_id",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "52sylu6udk1kc6r",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": null,
				"displayFields": null
			}
		}`), del_post_id); err != nil {
			return err
		}
		collection.Schema.AddField(del_post_id)

		return dao.SaveCollection(collection)
	})
}

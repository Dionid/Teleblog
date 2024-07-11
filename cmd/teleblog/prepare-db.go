package main

import (
	"strings"

	"github.com/Dionid/teleblog/libs/teleblog"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func preSeedDB(app *pocketbase.PocketBase) {
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		tag := teleblog.Tag{}
		err := teleblog.TagQuery(app.Dao()).
			Limit(1).
			One(&tag)
		if err != nil {
			return err
		}

		if tag.Id != "" {
			return nil
		}

		var posts []teleblog.Post

		err = teleblog.PostQuery(app.Dao()).
			All(&posts)
		if err != nil {
			return err
		}

		for _, post := range posts {
			tags, err := teleblog.ExtractTagsFromPost(post)
			if err != nil {
				return err
			}

			for _, tagValue := range tags {
				tag := teleblog.Tag{
					Value:  tagValue,
					PostId: post.Id,
					ChatId: post.ChatId,
				}

				err := app.Dao().Save(&tag)
				if err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") {
						continue
					}
					return err
				}
			}
		}

		return nil
	})
}

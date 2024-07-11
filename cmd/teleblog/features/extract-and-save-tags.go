package features

import (
	"fmt"
	"strings"

	"github.com/Dionid/teleblog/libs/teleblog"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

func ExtractAndSaveTags(app *pocketbase.PocketBase) error {
	_, err := app.Dao().DB().Delete((&teleblog.PostTag{}).TableName(), dbx.And()).Execute()
	if err != nil {
		return err
	}

	_, err = app.Dao().DB().Delete((&teleblog.Tag{}).TableName(), dbx.And()).Execute()
	if err != nil {
		return err
	}

	var posts []teleblog.Post

	err = teleblog.PostQuery(app.Dao()).
		OrderBy("created ASC").
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
				Value: tagValue,
			}

			if tagValue == "go" {
				fmt.Println("tagValue", tagValue, post.Id)
			}

			err := app.Dao().Save(&tag)
			if err != nil {
				if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
					return err
				}
				err = nil

				err = teleblog.TagQuery(app.Dao()).
					Where(dbx.HashExp{"value": tag.Value}).
					One(&tag)
				if err != nil {
					return err
				}
			}

			postTag := teleblog.PostTag{
				TagId:  tag.Id,
				PostId: post.Id,
				ChatId: post.ChatId,
			}

			err = app.Dao().Save(&postTag)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE constraint failed") {
					continue
				}
				return err
			}
		}
	}

	return nil
}

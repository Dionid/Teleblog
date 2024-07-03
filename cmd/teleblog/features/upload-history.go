package features

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Dionid/teleblog/libs/teleblog"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func UploadHistory(app core.App, history teleblog.History) error {
	chatId, err := history.AddMessageGetChatTgId()
	if err != nil {
		return err
	}

	var chat teleblog.Chat

	err = teleblog.ChatQuery(app.Dao()).Where(
		dbx.HashExp{"tg_chat_id": chatId},
	).Limit(1).One(&chat)
	if err != nil {
		return err
	}

	if chat.TgType != "channel" {
		return nil
	}

	var preparedPosts []teleblog.Post

	for _, message := range history.Messages {
		// # Check if exists
		total := struct {
			Total int64 `db:"total"`
		}{}

		err = teleblog.PostQuery(app.Dao()).
			Where(
				dbx.HashExp{"tg_post_id": message.Id, "chat_id": chat.Id},
			).
			Select(
				"count(*) as total",
			).
			One(&total)
		if err != nil {
			return err
		}

		if total.Total > 0 {
			continue
		}

		// # Remove files for now
		if message.File != nil {
			continue
		}

		// # Create new
		post := teleblog.Post{
			ChatId:             chat.Id,
			IsTgMessage:        true,
			IsTgHistoryMessage: true,
			Text:               "",
			TgMessageId:        message.Id,
		}

		// # post.Created
		if message.DateUnix != "" {
			i, err := strconv.ParseInt(message.DateUnix, 10, 64)
			if err != nil {
				return err
			}
			tm := time.Unix(i, 0)
			post.Created.Scan(tm)
		}

		// # post.TgHistoryEntities
		jsonEntities, err := json.Marshal(message.TextEntities)
		if err != nil {
			return err
		}

		err = post.TgHistoryEntities.Scan(jsonEntities)
		if err != nil {
			return err
		}

		// # post.TgMessageRaw
		jsonMessageRaw, err := json.Marshal(message)
		if err != nil {
			return err
		}

		err = post.TgMessageRaw.Scan(jsonMessageRaw)
		if err != nil {
			return err
		}

		preparedPosts = append(preparedPosts, post)
	}

	// # Save
	if len(preparedPosts) == 0 {
		return nil
	}

	for _, post := range preparedPosts {
		err = app.Dao().Save(&post)
		if err != nil {
			return err
		}
	}

	return nil
}

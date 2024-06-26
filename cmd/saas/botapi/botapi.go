package botapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dionid/teleadmin/libs/teleblog"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const ADD_CHANNEL_COMMAND_NAME = "addchannel"
const VERIFY_TOKEN_COMMAND_NAME = "verifytoken"

func InitBotCommands(b *telebot.Bot, app *pocketbase.PocketBase) {
	err := b.SetCommands([]telebot.Command{
		{Text: "start", Description: "start the bot"},
		{Text: VERIFY_TOKEN_COMMAND_NAME, Description: "send token to bind bot to your telebot account (e.g. /verifytoken YOUR_TOKEN)"},
		{Text: ADD_CHANNEL_COMMAND_NAME, Description: "send channel to create blog from it (e.g. /addchannel @YOUR_CHANNEL_NAME)"},
	})
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/start", func(c telebot.Context) error {
		return c.Reply("Hello! This is teleblog bot. Add it to your channel and get posts in your blog.")
	})

	b.Use(middleware.Recover(func(err error, ctx telebot.Context) {
		app.Logger().Error("Error in bot", err)
	}))

	VerifyTokenCommand(b, app)
	AddChannelCommand(b, app)

	b.Handle(telebot.OnChannelPost, func(c telebot.Context) error {
		chat := &teleblog.Chat{}

		err = teleblog.ChatQuery(app.Dao()).
			AndWhere(dbx.HashExp{"tg_chat_id": c.Chat().ID}).
			Limit(1).
			One(chat)
		if err != nil {
			return err
		}

		newPost := &teleblog.Post{
			ChatId:      chat.Id,
			IsTgMessage: true,
			Text:        c.Message().Text,
			TgMessageId: c.Message().ID,
		}

		jsonEntities, err := json.Marshal(c.Message().Entities)
		if err != nil {
			return err
		}

		err = newPost.TgEntities.Scan(jsonEntities)
		if err != nil {
			return err
		}

		jsonMessageRaw, err := json.Marshal(c.Message())
		if err != nil {
			return err
		}

		err = newPost.TgMessageRaw.Scan(jsonMessageRaw)
		if err != nil {
			return err
		}

		err = app.Dao().Save(newPost)
		if err != nil {
			return err
		}

		return nil
	})

	// # Created messages in channels, groups and bot
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		// # Check if it is from the channels chat
		// ...

		// # CHANNEL MESSAGES ARE ALSO HERE
		// ...

		app.Logger().Info("telebot.OnText")

		// # 0 if reply to something, or Post.Id if reply to post
		if c.Message().ReplyTo != nil {
			fmt.Println("c.Message().ReplyTo.OriginalMessageID", c.Message().ReplyTo.OriginalMessageID)
		}

		chat := &teleblog.Chat{}
		err = teleblog.ChatQuery(app.Dao()).
			AndWhere(dbx.HashExp{"tg_chat_id": c.Chat().ID}).
			Limit(1).
			One(chat)
		if err != nil {
			return err
		}

		if c.Message().FromChannel() {
			fmt.Println("Channel!", c.Message().Text)
		} else if c.Message().FromGroup() {
			fmt.Println("Group!", c.Message().Text)

			// # Forward from channel to group
			if c.Message().OriginalChat != nil && c.Message().OriginalChat.ID == chat.TgLinkedChatId {
				_, err := app.DB().Update(
					(&teleblog.Post{}).TableName(),
					map[string]interface{}{
						"tg_group_message_id": c.Message().ID,
					},
					dbx.HashExp{"tg_post_id": c.Message().OriginalMessageID},
				).Execute()

				return err
			}

			// # Bind by thread id
			if c.Message().ThreadID > 0 {
				post := &teleblog.Post{}

				err := teleblog.PostQuery(app.Dao()).
					AndWhere(dbx.HashExp{"tg_group_message_id": c.Message().ThreadID}).
					Limit(1).
					One(post)
				if err != nil {
					return err
				}

				newComment := &teleblog.Comment{
					ChatId:      chat.Id,
					PostId:      post.Id,
					Text:        c.Message().Text,
					TgMessageId: c.Message().ID,
				}

				if c.Message().ReplyTo != nil {
					newComment.TgReplyToMessageId = c.Message().ReplyTo.ID
				}

				jsonEntities, err := json.Marshal(c.Message().Entities)
				if err != nil {
					return err
				}

				err = newComment.TgEntities.Scan(jsonEntities)
				if err != nil {
					return err
				}

				jsonMessageRaw, err := json.Marshal(c.Message())
				if err != nil {
					return err
				}

				err = newComment.TgMessageRaw.Scan(jsonMessageRaw)
				if err != nil {
					return err
				}

				err = app.Dao().Save(newComment)
				if err != nil {
					return err
				}
			}
		} else {
			fmt.Println("Unknown", c.Message().Text)
		}

		return nil
	})

	// # Edited messages in channels and groups
	b.Handle(telebot.OnEditedChannelPost, func(c telebot.Context) error {
		fmt.Println("OnEditedChannelPost")

		chat := &teleblog.Chat{}

		err = teleblog.ChatQuery(app.Dao()).
			AndWhere(dbx.HashExp{"tg_chat_id": c.Chat().ID}).
			Limit(1).
			One(chat)
		if err != nil {
			return err
		}

		_, err := app.DB().Update(
			(&teleblog.Post{}).TableName(),
			map[string]interface{}{
				"text": c.Message().Text,
			},
			dbx.HashExp{"chat_id": chat.Id, "tg_post_id": c.Message().ID},
		).Execute()

		return err
	})

	b.Handle(telebot.OnEdited, func(c telebot.Context) error {
		fmt.Println("OnEdited", c.Message().Text)
		fmt.Println("c.Sender().ID", c.Sender().ID)

		chat := &teleblog.Chat{}

		err = teleblog.ChatQuery(app.Dao()).
			AndWhere(dbx.HashExp{"tg_chat_id": c.Chat().ID}).
			Limit(1).
			One(chat)
		if err != nil {
			return err
		}

		if c.Message().OriginalChat != nil && c.Message().OriginalChat.ID == chat.TgLinkedChatId {
			fmt.Println("FROM POST EDIT")
			return nil
		}

		_, err := app.DB().Update(
			(&teleblog.Comment{}).TableName(),
			map[string]interface{}{
				"text": c.Message().Text,
			},
			dbx.HashExp{"chat_id": chat.Id, "tg_comment_id": c.Message().ID},
		).Execute()

		return err
	})
}

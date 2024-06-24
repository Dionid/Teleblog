package botapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dionid/teleadmin/libs/teleblog"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"gopkg.in/telebot.v3"
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

		fmt.Println("OnChannelPost", c.Message().Entities)

		newPost := &teleblog.Post{
			ChatId:      chat.Id,
			IsTgMessage: true,
			Text:        c.Message().Text,
			TgPostId:    c.Message().ID,
		}

		jsonEntities, err := json.Marshal(c.Message().Entities)
		if err != nil {
			return err
		}

		err = newPost.TgEntities.Scan(jsonEntities)
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
			if c.Message().OriginalChat != nil && c.Message().OriginalChat.ID == chat.TgLinkedChatId {
				fmt.Println("This is from original", c.Message().Text)
				return nil
			}
		} else {
			fmt.Println("Bot!", c.Message().Text)
		}

		// fmt.Println("OnText", c.Message().FromChannel())
		// fmt.Println("OnText text", c.Message().Text)
		// fmt.Println("OnText c.Sender().ID", c.Sender().ID)

		return nil
	})

	// # Edited messages in channels and groups
	b.Handle(telebot.OnEdited, func(c telebot.Context) error {
		fmt.Println("OnEdited", c.Message().Text)
		fmt.Println("c.Sender().ID", c.Sender().ID)

		return nil
	})
}

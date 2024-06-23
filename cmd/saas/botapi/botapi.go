package botapi

import (
	"fmt"
	"log"

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

	// # Created messages in channels and groups
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		// # Check if it is from the channels chat
		// ...

		// # CHANNEL MESSAGES ARE ALSO HERE
		// ...

		fmt.Println("OnText", c.Message().Text)

		return nil
	})

	// # Edited messages in channels and groups
	b.Handle(telebot.OnEdited, func(c telebot.Context) error {
		fmt.Println("OnEdited", c.Message().Text)

		return nil
	})
}

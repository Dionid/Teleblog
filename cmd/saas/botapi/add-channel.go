package botapi

import (
	"fmt"

	"github.com/Dionid/teleadmin/libs/teleblog"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"gopkg.in/telebot.v3"
)

func AddChannelCommand(b *telebot.Bot, app *pocketbase.PocketBase) {
	b.Handle("/"+ADD_CHANNEL_COMMAND_NAME, func(c telebot.Context) error {
		user := &teleblog.User{}

		err := teleblog.UserQuery(app.Dao()).
			AndWhere(dbx.HashExp{"telegram_user_id": c.Sender().ID}).
			Limit(1).
			One(user)

		if err != nil {
			return c.Reply("You are not verified.")
		}

		tags := c.Args() // list of arguments splitted by a space

		if len(tags) == 0 {
			return c.Reply(fmt.Sprintf("You must provide channel name (e.g. /%s @YOUR_CHANNEL_NAME).", ADD_CHANNEL_COMMAND_NAME))
		}

		if len(tags) > 1 {
			return c.Reply("You can add only 1 channel at a time.")
		}

		channelUsername := tags[0]

		channel, err := b.ChatByUsername(channelUsername)
		if err != nil {
			return c.Reply("No channel like this found.")
		}

		if channel.Type != telebot.ChatChannel && channel.Type != telebot.ChatChannelPrivate {
			return c.Reply("This is not a channel.")
		}

		channelMember, err := b.ChatMemberOf(channel, c.Sender())
		if err != nil {
			return c.Reply("You cant add not your channels.")
		}

		if channelMember.Role != telebot.Administrator && channelMember.Role != telebot.Creator {
			return c.Reply("You are not the administrator of the channel.")
		}

		// TODO: # Check that channel.ID + user.Id is unique
		// ...

		newChannel := teleblog.Source{
			UserId:         user.Id,
			LinkedSourceId: "",

			Username:     channelUsername,
			ChatId:       channel.ID,
			Type:         string(channel.Type),
			LinkedChatId: channel.LinkedChatID,
		}

		if err := app.Dao().Save(&newChannel); err != nil {
			return err
		}

		// # Add linked chat
		linkedGroup, err := b.ChatByID(channel.LinkedChatID)
		if err != nil {
			return c.Reply("No channel chat like this found.")
		}

		channelGroupMember, err := b.ChatMemberOf(linkedGroup, c.Sender())
		if err != nil {
			return c.Reply("You cant add not your channels.")
		}

		fmt.Printf("channelGroupMember.Role: %s\n", channelGroupMember.Role)

		// TODO: # Check that channel.ID + user.Id is unique
		// ...

		newChannelGroup := teleblog.Source{
			UserId:         user.Id,
			LinkedSourceId: newChannel.Id,

			Username:     linkedGroup.Username,
			ChatId:       linkedGroup.ID,
			Type:         string(linkedGroup.Type),
			LinkedChatId: linkedGroup.LinkedChatID,
		}

		if err := app.Dao().Save(&newChannelGroup); err != nil {
			return err
		}

		return c.Reply("Channel and linked group are added.")
	})
}

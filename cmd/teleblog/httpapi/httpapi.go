package httpapi

import (
	"context"
	"embed"
	"encoding/json"
	"log"
	"os"

	"github.com/Dionid/teleblog/cmd/teleblog/httpapi/views"
	"github.com/Dionid/teleblog/libs/file"
	"github.com/Dionid/teleblog/libs/teleblog"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"gopkg.in/telebot.v3"
)

type Config struct {
	Env    string
	UserId string
}

//go:embed public
var publicAssets embed.FS

func InitApi(config Config, app core.App, gctx context.Context) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(apis.ActivityLogger(app))

		// # Static
		if config.Env == "PRODUCTION" {
			os.RemoveAll("./public")
			file.CopyFromEmbed(publicAssets, "public", "./public")
			e.Router.Static("/public", "./public")
		} else if config.Env == "LOCAL" {
			e.Router.Static("/public", "./httpapi/public")
		} else {
			log.Fatalf("Unknown env: %s", config.Env)
		}

		IndexPageHandler(config, e, app)

		e.Router.GET("/post/:id", func(c echo.Context) error {
			id := c.PathParam("id")

			post := views.PostPagePost{}

			err := teleblog.PostQuery(app.Dao()).Where(
				dbx.HashExp{"id": id},
			).Limit(1).One(&post)
			if err != nil {
				return err
			}

			jb, err := post.Post.TgMessageRaw.MarshalJSON()
			if err != nil {
				return err
			}

			if post.IsTgHistoryMessage {
				rawMessage := teleblog.HistoryMessage{}

				err = json.Unmarshal(jb, &rawMessage)
				if err != nil {
					return err
				}

				post.TextWithMarkup = teleblog.FormHistoryTextWithMarkup(rawMessage.TextEntities)
			} else {
				rawMessage := telebot.Message{}

				err = json.Unmarshal(jb, &rawMessage)
				if err != nil {
					return err
				}

				post.TextWithMarkup, err = teleblog.FormWebhookTextMarkup(post.Text, rawMessage.Entities)
				if err != nil {
					return err
				}
			}

			chat := teleblog.Chat{}

			err = teleblog.ChatQuery(app.Dao()).Where(
				dbx.HashExp{"id": post.ChatId},
			).Limit(1).One(&chat)
			if err != nil {
				return err
			}

			comments := []*views.PostPageComment{}

			err = teleblog.CommentQuery(app.Dao()).Where(
				dbx.HashExp{"post_id": id},
			).All(&comments)
			if err != nil {
				return err
			}

			for _, comment := range comments {
				jb, err := comment.TgMessageRaw.MarshalJSON()
				if err != nil {
					return err
				}

				if comment.IsTgHistoryMessage {
					rawMessage := teleblog.HistoryMessage{}

					err = json.Unmarshal(jb, &rawMessage)
					if err != nil {
						return err
					}

					comment.AuthorTitle = rawMessage.From

					comment.TextWithMarkup = teleblog.FormHistoryTextWithMarkup(rawMessage.TextEntities)
				} else {
					rawMessage := telebot.Message{}

					err = json.Unmarshal(jb, &rawMessage)
					if err != nil {
						return err
					}

					if rawMessage.Sender.IsBot && rawMessage.SenderChat != nil {
						comment.AuthorTitle = rawMessage.SenderChat.Title
						comment.AuthorUsername = &rawMessage.SenderChat.Username
					} else {
						comment.AuthorTitle = rawMessage.Sender.FirstName + " " + rawMessage.Sender.LastName
						comment.AuthorUsername = &rawMessage.Sender.Username
					}

					comment.TextWithMarkup, err = teleblog.FormWebhookTextMarkup(comment.Text, rawMessage.Entities)
					if err != nil {
						return err
					}
				}
			}

			component := views.PostPage(chat, post, comments)

			return component.Render(c.Request().Context(), c.Response().Writer)
		})

		return nil
	})
}

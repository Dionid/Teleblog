package httpapi

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Dionid/teleadmin/cmd/saas/httpapi/views"
	"github.com/Dionid/teleadmin/libs/file"
	"github.com/Dionid/teleadmin/libs/teleblog"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type Config struct {
	Env string
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

		type PostPageFilters struct {
			Page    int64 `query:"page"`
			PerPage int64 `query:"per_page"`
		}

		e.Router.GET("", func(c echo.Context) error {
			chats := []teleblog.Chat{}

			err := teleblog.ChatQuery(app.Dao()).Where(
				dbx.HashExp{"user_id": "wujghdqjma1fydw", "tg_type": "channel"},
			).All(&chats)
			if err != nil {
				return err
			}

			if len(chats) == 0 {
				// TODO: Change
				return c.JSON(200, []teleblog.Post{})
			}

			chatIds := []interface{}{}
			for _, chat := range chats {
				chatIds = append(chatIds, chat.Id)
			}

			// # Filters
			var filters PostPageFilters

			if err := c.Bind(&filters); err != nil {
				return err
			}

			// Query
			baseQuery := teleblog.PostQuery(app.Dao()).
				Where(
					dbx.In("post.chat_id", chatIds...),
				)

			// ## Total
			total := struct {
				Total int64 `db:"total"`
			}{}

			err = baseQuery.Select(
				"count(post.id) as total",
			).One(&total)
			if err != nil {
				return err
			}

			fmt.Println("total", total)

			// ## Posts
			posts := []*views.InpexPagePost{}
			contentQuery := baseQuery.Select(
				"post.*",
				"count(comment.id) as comments_count",
				"chat.tg_username as tg_chat_username",
			).
				LeftJoin(
					"comment",
					dbx.NewExp("comment.post_id = post.id"),
				).
				LeftJoin(
					"chat",
					dbx.NewExp("chat.id = post.chat_id"),
				).
				GroupBy("post.id").
				OrderBy("post.created desc")

			// ## Filters
			// ### Per page
			perPage := filters.PerPage

			if perPage == 0 {
				perPage = 10
			} else if perPage > 100 {
				perPage = 100
			}

			contentQuery = contentQuery.Limit(perPage)

			// ### Current page
			currentPage := filters.Page
			if currentPage == 0 {
				currentPage = 1
			}

			contentQuery = contentQuery.Offset((currentPage - 1) * perPage)

			err = contentQuery.
				All(&posts)
			if err != nil {
				return err
			}

			for _, post := range posts {
				markup, err := teleblog.AddMarkupToText(post.Text, post.TgEntities)
				if err != nil {
					return err
				}

				fmt.Println("markup", markup)
				post.TextWithMarkup = markup
			}

			pagination := views.PaginationData{
				Total:       total.Total,
				PerPage:     perPage,
				CurrentPage: currentPage,
			}

			component := views.IndexPage(pagination, posts)

			return component.Render(c.Request().Context(), c.Response().Writer)
		})

		e.Router.GET("/post/:id", func(c echo.Context) error {
			id := c.PathParam("id")

			post := views.PostPagePost{}

			err := teleblog.PostQuery(app.Dao()).Where(
				dbx.HashExp{"id": id},
			).Limit(1).One(&post)
			if err != nil {
				return err
			}

			post.TextWithMarkup, err = teleblog.AddMarkupToText(post.Text, post.TgEntities)
			if err != nil {
				return err
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
				rawMessage := struct {
					From struct {
						IsBot     bool   `json:"is_bot"`
						Username  string `json:"username"`
						FirstName string `json:"first_name"`
						LastName  string `json:"last_name"`
					} `json:"from"`
					SenderChat struct {
						Title    string `json:"title"`
						Username string `json:"username"`
					} `json:"sender_chat"`
				}{}

				jb, err := comment.TgMessageRaw.MarshalJSON()
				if err != nil {
					return err
				}

				err = json.Unmarshal(jb, &rawMessage)
				if err != nil {
					return err
				}

				fmt.Println("rawMessage", rawMessage.SenderChat.Title)

				if rawMessage.From.IsBot {
					comment.AuthorTitle = rawMessage.SenderChat.Title
					comment.AuthorUsername = rawMessage.SenderChat.Username
				} else {
					comment.AuthorTitle = rawMessage.From.FirstName + " " + rawMessage.From.LastName
					comment.AuthorUsername = rawMessage.From.Username
				}
			}

			component := views.PostPage(chat, post, comments)

			return component.Render(c.Request().Context(), c.Response().Writer)
		})

		return nil
	})
}

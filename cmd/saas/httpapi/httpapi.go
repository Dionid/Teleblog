package httpapi

import (
	"context"

	"github.com/Dionid/teleadmin/cmd/saas/httpapi/views"
	"github.com/Dionid/teleadmin/libs/teleblog"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func InitApi(app core.App, gctx context.Context) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(apis.ActivityLogger(app))

		e.Router.GET("", func(c echo.Context) error {
			chats := []teleblog.Chat{}

			err := teleblog.ChatQuery(app.Dao()).Where(
				dbx.HashExp{"user_id": "wujghdqjma1fydw", "tg_type": "channel"},
			).All(&chats)
			if err != nil {
				return err
			}

			if len(chats) == 0 {
				return c.JSON(200, []teleblog.Post{})
			}

			chatIds := []interface{}{}
			for _, chat := range chats {
				chatIds = append(chatIds, chat.Id)
			}

			posts := []teleblog.Post{}

			err = teleblog.PostQuery(app.Dao()).Where(
				dbx.In("chat_id", chatIds...),
			).All(&posts)
			if err != nil {
				return err
			}

			component := views.IndexPage(posts)

			return component.Render(c.Request().Context(), c.Response().Writer)
		})

		e.Router.GET("/post/:id", func(c echo.Context) error {
			id := c.PathParam("id")

			post := teleblog.Post{}

			err := teleblog.PostQuery(app.Dao()).Where(
				dbx.HashExp{"id": id},
			).Limit(1).One(&post)
			if err != nil {
				return err
			}

			comments := []teleblog.Comment{}

			err = teleblog.CommentQuery(app.Dao()).Where(
				dbx.HashExp{"post_id": id},
			).All(&comments)
			if err != nil {
				return err
			}

			component := views.PostPage(post, comments)

			return component.Render(c.Request().Context(), c.Response().Writer)
		})

		return nil
	})
}

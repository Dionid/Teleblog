package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Dionid/teleblog/cmd/teleblog/botapi"
	"github.com/Dionid/teleblog/cmd/teleblog/features"
	"github.com/Dionid/teleblog/cmd/teleblog/httpapi"
	_ "github.com/Dionid/teleblog/cmd/teleblog/pb_migrations"
	"github.com/Dionid/teleblog/libs/teleblog"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/spf13/cobra"
	"gopkg.in/telebot.v3"
)

func main() {
	config, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	gctx, _ := context.WithCancel(context.Background())

	// # Pocketbase
	app := pocketbase.New()

	// # Send verification email on sign-up
	app.OnRecordAfterCreateRequest("users").Add(func(e *core.RecordCreateEvent) error {
		return mails.SendRecordVerification(app, e.Record)
	})

	app.RootCmd.AddCommand(&cobra.Command{
		Use: "upload-history",
		Run: func(cmd *cobra.Command, args []string) {
			defer (func() {
				if r := recover(); r != nil {
					log.Fatal("recover", r)
				}
			})()

			fileName := "result.json"

			if len(args) > 0 {
				fileName = args[0]
			}

			file, err := os.ReadFile(fileName)
			if err != nil {
				log.Fatal(err)
			}

			var history teleblog.History
			err = json.Unmarshal(file, &history)
			if err != nil {
				log.Fatal(err)
			}

			err = features.UploadHistory(app, history)
			if err != nil {
				log.Fatal(err)
			}

			app.Logger().Info("Done")
		},
	})

	// # Migrations
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	curPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
		Dir:         path.Join(curPath, "pb_migrations"),
	})

	// # Telegram
	pref := telebot.Settings{
		Verbose: true,
		Token:   config.TelegramNotToken,
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second, AllowedUpdates: telebot.AllowedUpdates},
		OnError: func(err error, c telebot.Context) {
			app.Logger().Error("Error in bot", "error:", err)
		},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	botapi.InitBotCommands(b, app)

	// # HTTP API

	httpapi.InitApi(httpapi.Config{
		Env:    config.Env,
		UserId: config.UserId,
	}, app, gctx)

	// # Start bot
	go b.Start()

	// # Start app
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

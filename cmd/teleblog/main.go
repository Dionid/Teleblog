package main

import (
	"context"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Dionid/teleadmin/cmd/teleblog/botapi"
	"github.com/Dionid/teleadmin/cmd/teleblog/httpapi"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
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
		Env: config.Env,
	}, app, gctx)

	// # Send verification email on sign-up
	app.OnRecordAfterCreateRequest("users").Add(func(e *core.RecordCreateEvent) error {
		return mails.SendRecordVerification(app, e.Record)
	})

	// # Start bot
	go b.Start()

	// # Start app
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
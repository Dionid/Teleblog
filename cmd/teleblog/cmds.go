package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Dionid/teleblog/cmd/teleblog/features"
	"github.com/Dionid/teleblog/libs/teleblog"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func AdditionalCommands(app *pocketbase.PocketBase) {
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

	app.RootCmd.AddCommand(&cobra.Command{
		Use: "extract-tags",
		Run: func(cmd *cobra.Command, args []string) {
			defer (func() {
				if r := recover(); r != nil {
					log.Fatal("recover", r)
				}
			})()

			err := features.ExtractAndSaveTags(app)
			if err != nil {
				log.Fatal(err)
			}

			app.Logger().Info("Done")
		},
	})
}

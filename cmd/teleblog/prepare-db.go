package main

import (
	"strings"

	"github.com/Dionid/teleblog/cmd/teleblog/features"
	"github.com/Dionid/teleblog/libs/teleblog"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func preSeedDB(app *pocketbase.PocketBase) {
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		var existingTags []teleblog.Tag
		err := teleblog.TagQuery(app.Dao()).
			Limit(1).
			All(&existingTags)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil
			}
			return err
		}

		if len(existingTags) > 0 {
			return nil
		}

		return features.ExtractAndSaveAllTags(app)
	})
}

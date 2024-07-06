package templu

import (
	"context"
	"os"
)

func PathWithVersion(ctx context.Context, url string) string {
	if ctxAppVersion, ok := ctx.Value("APP_VERSION").(string); ok && ctxAppVersion != "" {
		return url + "?v=" + ctxAppVersion
	}

	appVersion := os.Getenv("APP_VERSION")
	if appVersion != "" {
		return url + "?v=" + os.Getenv("APP_VERSION")
	}

	return url + "?v=0.0.1"
}

package main

import (
	"log/slog"
	"os"

	"github.com/josestg/swe-be-mono/internal/app"
	"github.com/josestg/swe-be-mono/internal/app/adminrestful"
)

// These variables are set by the build process.
// see: https://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning/11355611#11355611.
var (
	buildName    = "unset"
	buildTime    = "unset"
	buildVersion = "unset"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	slog.SetDefault(log)

	info, err := app.NewInfo(buildName, buildTime, buildVersion)
	if err != nil {
		log.Warn("failed to create app info", "error", err)
	}

	if err := adminrestful.Run(log, info); err != nil {
		log.Error("run failed", "error", err)
		os.Exit(1)
	}
}

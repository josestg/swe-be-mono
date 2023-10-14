package main

import (
	"log/slog"
	"os"

	"github.com/josestg/swe-be-mono/internal/app"
	"github.com/josestg/swe-be-mono/internal/app/adminrestful"
	"github.com/josestg/swe-be-mono/internal/config"
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

	cfg, err := config.New(buildName, buildTime, buildVersion)
	if err != nil {
		log.Error("failed to create app info", "error", err)
		os.Exit(1)
	}

	if err := app.Run(log, cfg, adminrestful.AppFactory); err != nil {
		log.Error("app run failed", "error", err)
		os.Exit(1)
	}
}

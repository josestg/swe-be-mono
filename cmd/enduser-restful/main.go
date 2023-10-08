package main

import (
	"fmt"
	"os"

	"github.com/josestg/swe-be-mono/internal/app"
	"github.com/josestg/swe-be-mono/internal/app/enduserrestful"
)

// These variables are set by the build process.
// see: https://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning/11355611#11355611.
var (
	buildName    = "unset"
	buildTime    = "unset"
	buildVersion = "unset"
)

func main() {
	info := app.Info{
		Name:         buildName,
		BuildTime:    buildTime,
		BuildVersion: buildVersion,
	}

	if err := enduserrestful.Run(info); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[%s]: run failed, got error: %v\n", info.Name, err)
		os.Exit(1)
	}
}

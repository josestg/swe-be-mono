package main

import (
	"fmt"
	"github.com/josestg/swe-be-mono/internal/app"
	"github.com/josestg/swe-be-mono/internal/app/enduserrestful"
	"os"
)

func main() {
	info := app.Info{
		Name:         "enduser-restful",
		BuildTime:    "TODO",
		BuildVersion: "TODO",
	}

	if err := enduserrestful.Run(info); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[%s]: run failed, got error: %v\n", info.Name, err)
		os.Exit(1)
	}
}

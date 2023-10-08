package main

import (
	"fmt"
	"github.com/josestg/swe-be-mono/internal/app"
	"github.com/josestg/swe-be-mono/internal/app/adminrestful"
	"os"
)

func main() {
	info := app.Info{
		Name:         "admin-restful",
		BuildTime:    "TODO",
		BuildVersion: "TODO",
	}

	if err := adminrestful.Run(info); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[%s]: run failed, got error: %v\n", info.Name, err)
		os.Exit(1)
	}
}

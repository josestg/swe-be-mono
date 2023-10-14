package enduserrestful

import (
	"log/slog"

	"github.com/josestg/swe-be-mono/internal/app"
)

// Run is the entrypoint of the enduser-restful application.
func Run(log *slog.Logger, app app.Info) error {
	log.Info("app started", "app", app)
	defer log.Info("app stopped", "app", app)
	return nil
}

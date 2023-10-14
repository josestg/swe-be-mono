package adminrestful

import (
	"log/slog"
	"net/http"

	"github.com/josestg/swe-be-mono/internal/app"
	"github.com/josestg/swe-be-mono/pkg/httpkit"
)

// App is the admin-restful application.
type App struct {
	log *slog.Logger
}

// AppFactory is the factory for creating the admin-restful application.
func AppFactory(log *slog.Logger) app.App {
	return &App{
		log: log,
	}
}

// DocHandler returns the handler for the admin-restful documentation.
func (a *App) DocHandler() http.Handler { return _docHandler }

// APIHandler returns the handler for the admin-restful APIs.
func (a *App) APIHandler() http.Handler {
	mux := httpkit.NewServeMux()
	return mux
}

// _docHandler is the default handler for docs endpoint.
var _docHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	msg := `{"message": "please run with 'swagger_docs_enabled' build tag to enable swagger docs"}`
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, msg, http.StatusNotImplemented)
})

// SetDocHandler sets the handler for docs endpoint.
// lint:ignore U1000 because this is used by build tag.
func SetDocHandler(handler http.Handler) { _docHandler = handler }

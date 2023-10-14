package httphandler

import (
	"net/http"

	"github.com/josestg/swe-be-mono/internal/config"

	"github.com/josestg/swe-be-mono/pkg/httpkit"
)

// System is a handler for serving system information and health checks.
type System struct {
	app config.AppInfo
}

func ServeSystem(mux *httpkit.ServeMux, app config.AppInfo) {
	sys := &System{app: app}
	mux.Route(sys.Info())
	mux.Route(sys.Health())
}

func (h *System) MountTo(mux *httpkit.ServeMux) {
	mux.Route(h.Info())
	mux.Route(h.Health())
}

// Info returns the application information.
func (h *System) Info() httpkit.Route {
	return httpkit.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/sys/info",
		Handler: h.info,
	}
}

// Health returns the application health status.
func (h *System) Health() httpkit.Route {
	return httpkit.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/sys/health",
		Handler: h.health,
	}
}

func (h *System) info(w http.ResponseWriter, _ *http.Request) error {
	return httpkit.WriteJSON(w, h.app, http.StatusOK)
}

func (h *System) health(w http.ResponseWriter, _ *http.Request) error {
	dependencies := []struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}{
		{
			Name:   "database",
			Status: "healthy",
		},
	}

	return httpkit.WriteJSON(w, dependencies, http.StatusOK)
}

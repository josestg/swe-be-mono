package httphandler

import (
	"net/http"

	"github.com/josestg/swe-be-mono/internal/config"
	"github.com/josestg/swe-be-mono/internal/domain/system"
	"github.com/josestg/swe-be-mono/internal/kernel"
	"github.com/josestg/swe-be-mono/pkg/httpkit"
)

// System is a handler for serving system information and health checks.
type System struct {
	app config.AppInfo
}

// ServeSystem registers the system handler to the given mux.
func ServeSystem(mux *httpkit.ServeMux, app config.AppInfo) {
	sys := &System{app: app}
	mux.Route(sys.Info())
	mux.Route(sys.Health())
}

// Info returns the application information.
//
//	@Tags			System
//	@Summary		Application information.
//	@Description	Returns the application information.
//	@Produce		json
//	@Success		200	{object}	kernel.HttpRes[config.AppInfo]
//	@Router			/system/info [get]
func (h *System) Info() httpkit.Route {
	return httpkit.Route{
		Method:  http.MethodGet,
		Path:    "/system/info",
		Handler: h.info,
	}
}

// Health returns the application health status.
//
//	@Tags			System
//	@Summary		Application health.
//	@Description	Returns the application health status.
//	@Produce		json
//	@Success		200	{object}	kernel.HttpRes[[]system.HealthRes]
//	@Router			/system/health [get]
func (h *System) Health() httpkit.Route {
	return httpkit.Route{
		Method:  http.MethodGet,
		Path:    "/system/health",
		Handler: h.health,
	}
}

func (h *System) info(w http.ResponseWriter, _ *http.Request) error {
	res := kernel.NewHttpResBuilder(h.app).Build()
	return httpkit.WriteJSON(w, res, res.Code)
}

func (h *System) health(w http.ResponseWriter, _ *http.Request) error {
	dependencies := []system.HealthRes{
		{
			Name:   "HTTP Server",
			Status: system.StatusHealthy,
		},
		{
			Name:   "MySQL",
			Status: system.StatusUnhealthy,
		},
	}

	res := kernel.NewHttpResBuilder(dependencies).Build()
	return httpkit.WriteJSON(w, res, res.Code)
}

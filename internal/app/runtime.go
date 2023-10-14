package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"syscall"

	"github.com/josestg/swe-be-mono/internal/config"

	"github.com/josestg/swe-be-mono/internal/httphandler"

	"github.com/josestg/swe-be-mono/pkg/httpkit"
)

// Run is the entrypoint of the for the application.
func Run(log *slog.Logger, cfg *config.Config, factory Factory) error {
	log.Info("app started", "app", cfg.AppInfo)
	defer log.Info("app stopped", "app", cfg.AppInfo)

	router := newRouter(log, cfg, factory)
	return listenAndServe(log, cfg.HttpServer, router)
}

// newRouter returns the complete http.Handler for the application.
// Including the Application APIs, Documentation and System APIs.
func newRouter(log *slog.Logger, cfg *config.Config, factory Factory) http.Handler {
	app := factory.New(log)

	// mux in here is a root mux for splitting the traffic to different handlers based on the path prefix.
	mux := http.NewServeMux()
	mux.Handle("/docs/", app.DocHandler())
	mux.Handle("/api/v1/", app.APIHandler())
	mux.Handle("/systems/", systemHandler(cfg.AppInfo))
	return mux
}

// systemHandler is a handler for serving system information and health checks.
func systemHandler(info config.AppInfo) http.Handler {
	mux := httpkit.NewServeMux()
	httphandler.ServeSystem(mux, info)
	return mux
}

// listenAndServe starts the http server and gracefully shutdowns on signals received.
func listenAndServe(log *slog.Logger, cfg httpkit.RunConfig, mux http.Handler) error {
	srv := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  cfg.RequestReadTimeout,
		WriteTimeout: cfg.RequestWriteTimeout,
	}

	run := httpkit.NewGracefulRunner(&srv,
		httpkit.RunOpts.WaitTimeout(cfg.ShutdownTimeout),
		httpkit.RunOpts.Signals(syscall.SIGINT, syscall.SIGTERM),
		httpkit.RunOpts.EventListener(func(evt httpkit.RunEvent, data string) {
			log.Info("http graceful runtime", "event", evt, "data", data)
		}),
	)

	return run.ListenAndServe()
}
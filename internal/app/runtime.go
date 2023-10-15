package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"syscall"

	"github.com/josestg/swe-be-mono/internal/httpmiddleware"

	"github.com/josestg/swe-be-mono/internal/config"

	"github.com/josestg/swe-be-mono/internal/httphandler"

	"github.com/josestg/swe-be-mono/pkg/httpkit"
)

// Run is the entrypoint of the for the application.
func Run(log *slog.Logger, cfg *config.Config, factory Factory) error {
	log.Info("app started", "app", cfg.AppInfo)
	defer log.Info("app stopped", "app", cfg.AppInfo)

	router := newRouter(cfg, factory)
	return listenAndServe(log, cfg.HttpServer, router)
}

// newRouter returns the complete http.Handler for the application.
// Including the Application APIs, Documentation and System APIs.
func newRouter(cfg *config.Config, factory Factory) http.Handler {
	app := factory.New(cfg)

	// dynamically get the path prefix for the application.
	prefix := app.BasePath()

	// mux in here is a root mux for splitting the traffic to different handlers based on the path prefix.
	mux := http.NewServeMux()
	mux.Handle(prefix+"/docs/", app.DocHandler())
	mux.Handle(prefix+"/api/v1/", http.StripPrefix(prefix, app.APIHandler()))
	mux.Handle(prefix+"/system/", http.StripPrefix(prefix, systemHandler(cfg.AppInfo)))

	mid := httpkit.ReduceNetMiddleware(httpmiddleware.CORS(cfg.HttpCORS))
	return mid.Then(mux)
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
			switch evt {
			default:
				log.Info(data)
			case httpkit.RunEventAddr:
				log.Info("http server listening", "addr", data)
			case httpkit.RunEventSignal:
				log.Info("http server received shutdown signal", "signal", data)
			}
		}),
	)

	return run.ListenAndServe()
}

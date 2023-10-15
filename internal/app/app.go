package app

import (
	"net/http"

	"github.com/josestg/swe-be-mono/internal/config"
)

// App is contract for API application that can be run in application runtime.
type App interface {
	// APIHandler returns the handler for the Applications APIs.
	APIHandler() http.Handler

	// DocHandler returns the handler for the Applications Documentation.
	DocHandler() http.Handler

	// BasePath returns the base path for the application.
	BasePath() string
}

// Factory is a function that creates an instance of the application.
type Factory func(cfg *config.Config) App

// New is a syntactic sugar for applying the factory.
func (f Factory) New(cfg *config.Config) App { return f(cfg) }

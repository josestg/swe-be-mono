package app

import (
	"log/slog"
	"net/http"
)

// App is contract for API application that can be run in application runtime.
type App interface {
	// APIHandler returns the handler for the Applications APIs.
	APIHandler() http.Handler

	// DocHandler returns the handler for the Applications Documentation.
	DocHandler() http.Handler
}

// Factory is a function that creates an instance of the application.
type Factory func(log *slog.Logger) App

// New is a syntactic sugar for applying the factory.
func (f Factory) New(log *slog.Logger) App { return f(log) }

package httpmiddleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/josestg/problemdetail"
	"github.com/josestg/swe-be-mono/internal/business"
	"github.com/josestg/swe-be-mono/pkg/httpkit"
)

// LogAndErrHandling is a middleware that logs the request and response and
// handles error.
func LogAndErrHandling(log *slog.Logger) httpkit.MuxMiddleware {
	return func(next httpkit.Handler) httpkit.Handler {
		return httpkit.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			rec, ok := httpkit.GetLogEntry(w)
			if !ok {
				return fmt.Errorf("missing log entry: method=%s path=%s", r.Method, r.URL.Path)
			}

			err := next.ServeHTTP(w, r)
			if err == nil {
				log.LogAttrs(r.Context(), slog.LevelInfo, "completed",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("uri", r.RequestURI),
					slog.Int("status", rec.StatusCode),
					slog.Duration("latency", time.Duration(rec.RespondedAt-rec.RequestedAt)),
				)
				return nil
			}

			err = MapError(w, err)
			var resolvedErr *httpkit.ResolvedError
			if !errors.As(err, &resolvedErr) {
				log.LogAttrs(r.Context(), slog.LevelError, "unresolved_error",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("uri", r.RequestURI),
					slog.Int("status", rec.StatusCode),
					slog.Duration("latency", time.Duration(rec.RespondedAt-rec.RequestedAt)),
					slog.Any("error", err),
				)
			} else {
				log.LogAttrs(r.Context(), slog.LevelInfo, "resolved_error",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("uri", r.RequestURI),
					slog.Int("status", rec.StatusCode),
					slog.Duration("latency", time.Duration(rec.RespondedAt-rec.RequestedAt)),
					slog.Any("error", resolvedErr.Err),
				)
			}

			// As the error has been properly handled, we can return nil to
			// indicate to the subsequent chain that the error has been taken
			// care of.
			return nil
		})
	}
}

// MapError maps the error to an HTTP response and marks the error as resolved if
// it is successfully mapped.
func MapError(w http.ResponseWriter, err error) error {
	var pd problemdetail.ProblemDetailer
	if !errors.As(err, &pd) {
		// untyped error for generic error handling.
		untyped := problemdetail.New(
			problemdetail.Untyped,
			problemdetail.WithValidateLevel(problemdetail.LStandard),
		)
		return sendJSONError(w, http.StatusInternalServerError, untyped, err, false)
	}

	switch pd.Kind() {
	case business.PDTypeEmailAlreadyTaken:
		return sendJSONError(w, http.StatusConflict, pd, err, true)
	case business.PDTypeUserNotFound:
		return sendJSONError(w, http.StatusNotFound, pd, err, true)
	case business.PDTypeInvalidArguments:
		return sendJSONError(w, http.StatusBadRequest, pd, err, true)
	}

	return fmt.Errorf("could not map error: %w", err)
}

// sendJSONError sends the error as a JSON response.
func sendJSONError(w http.ResponseWriter, code int, data problemdetail.ProblemDetailer, err error, resolved bool) error {
	wErr := problemdetail.WriteJSON(w, data, code)
	if wErr != nil {
		// if we fail to write the error to the response,
		// we will encounter two errors: the error that needs to be handled
		// and the error that occurs when writing to the response.
		// Therefore, it is crucial to chain the errors to ensure that
		// we do not lose any of them.
		return errors.Join(wErr, err)
	}

	// if the error is resolved, we mark it as resolved so that the subsequent
	// middleware can determine whether the error has been handled or not.
	if resolved {
		return httpkit.ResolveError(err)
	}

	// we return the original error so that we can log it.
	return err
}

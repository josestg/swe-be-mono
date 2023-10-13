package httpkit

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// ReadJSON reads json from the reader and decodes it to the data.
// By default, it disallows unknown fields.
func ReadJSON(r io.Reader, data any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(data)
}

// WriteJSON writes the data to the response writer as JSON.
// By default, it sets the content type to application/json; charset=utf-8.
func WriteJSON(w http.ResponseWriter, data any, code int) error {
	writeContentTypeAndStatus(w, contentTypeApplicationJSONCharsetUTF8, code)
	return json.NewEncoder(w).Encode(data)
}

// WriteJSONProblemDetail writes the problem detail to the response writer as JSON.
// The content type is set to application/problem+json; charset=utf-8.
// The status code will be set to both ProblemDetail.Status and http.ResponseWriter.
//
// If the problem detail is invalid, an error is returned.
func WriteJSONProblemDetail(w http.ResponseWriter, pd ProblemDetailer, code int) error {
	if err := prepareProblemDetail(pd, code); err != nil {
		return fmt.Errorf("WriteJSONProblemDetail: %w", err)
	}
	writeContentTypeAndStatus(w, contentTypeJSONProblemDetail, code)
	return json.NewEncoder(w).Encode(pd)
}

// WriteXMLProblemDetail writes the problem detail to the response writer as XML.
// The content type is set to application/problem+xml; charset=utf-8.
// The status code will be set to both ProblemDetail.Status and http.ResponseWriter.
//
// If the problem detail is invalid, an error is returned.
func WriteXMLProblemDetail(w http.ResponseWriter, pd ProblemDetailer, code int) error {
	if err := prepareProblemDetail(pd, code); err != nil {
		return fmt.Errorf("WriteXMLProblemDetail: %w", err)
	}
	writeContentTypeAndStatus(w, contentTypeXMLProblemDetail, code)
	return xml.NewEncoder(w).Encode(pd)
}

func prepareProblemDetail(pd ProblemDetailer, code int) error {
	pd.WriteStatus(code)
	if err := pd.Validate(); err != nil {
		return fmt.Errorf("problem detail: validate: %w", err)
	}
	return nil
}

// writeContentTypeAndStatus writes the content type and status code to the response writer.
func writeContentTypeAndStatus(w http.ResponseWriter, value string, code int) {
	w.Header().Add("Content-Type", value)
	w.WriteHeader(code)
}

// MIMETypes
const (
	charsetUTF8                           = "charset=UTF-8"
	contextTypeApplicationJSON            = "application/json"
	contentTypeApplicationJSONCharsetUTF8 = contextTypeApplicationJSON + "; " + charsetUTF8
	contentTypeXMLProblemDetail           = "application/problem+xml" + "; " + charsetUTF8
	contentTypeJSONProblemDetail          = "application/problem+json" + "; " + charsetUTF8
)

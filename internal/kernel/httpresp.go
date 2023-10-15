package kernel

import (
	"net/http"
	"time"
)

// HttpRes is a base template for HTTP response.
// swagger:response kernel.HttpResp
type HttpRes[T any] struct {
	// Code is the http status code. The code must be in the range of 200-299,
	// for other codes, the response is considered as an error and please use the Problem Details.
	// Default value is 200.
	Code int `json:"code"`

	// Data is the response data.
	Data T `json:"data"`

	// Desc is a short and human-readable description of the response.
	// This field is optional.
	Desc string `json:"desc,omitempty"`

	// Time is the time in unix milliseconds that describes when the response is
	// created.
	// Default value is the current time.
	Time int64 `json:"time"`
} //@name kernel.HttpResp

// HttpResBuilder is a builder for HttpRes.
type HttpResBuilder[T any] struct {
	state HttpRes[T]
}

// NewHttpResBuilder creates a new HttpResBuilder with the given data and default
// status code (200) and time (current time).
func NewHttpResBuilder[T any](data T) *HttpResBuilder[T] {
	return &HttpResBuilder[T]{
		state: HttpRes[T]{
			Data: data,
			Code: http.StatusOK,
			Time: time.Now().UnixMilli(),
		},
	}
}

// Code sets the status code.
func (b *HttpResBuilder[T]) Code(code int) *HttpResBuilder[T] {
	b.state.Code = code
	return b
}

// Desc sets the description.
func (b *HttpResBuilder[T]) Desc(desc string) *HttpResBuilder[T] {
	b.state.Desc = desc
	return b
}

// Time sets the time.
func (b *HttpResBuilder[T]) Time(epochMillis int64) *HttpResBuilder[T] {
	b.state.Time = epochMillis
	return b
}

// Build returns the HttpRes that is built.
func (b *HttpResBuilder[T]) Build() HttpRes[T] { return b.state }

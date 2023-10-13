package kernel

type ProblemDetail struct {
	// Type is a URI reference [RFC3986] that identifies the problem type.
	Type string `json:"type"`

	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title"`

	// Status the HTTP status code.
	Status int `json:"status"`

	// Detail is a human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// Instance is a URI reference that identifies the specific occurrence of the problem.
	Instance string `json:"instance,omitempty"`
}

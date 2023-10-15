package httpmiddleware

import (
	"github.com/josestg/swe-be-mono/pkg/httpkit"
	"github.com/rs/cors"
)

// CORS is a middleware that handles CORS.
func CORS(cfg cors.Options) httpkit.NetMiddleware { return cors.New(cfg).Handler }

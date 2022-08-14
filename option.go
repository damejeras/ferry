package ferry

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

type Option func(mux *mux)

func WithOpenApiSpec(mod func(spec *openapi3.Spec)) Option {
	return func(s *mux) {
		s.openapiSpecMod = mod
	}
}

func WithErrorHandler(handler ErrorHandler) Option {
	return func(s *mux) {
		s.errHandler = handler
	}
}

func WithMiddleware(mw func(http.Handler) http.Handler) Option {
	return func(s *mux) {
		s.middleware = append(s.middleware, mw)
	}
}

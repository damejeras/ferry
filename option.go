package ferry

import (
	"strings"
)

type Option func(mux *ServeMux)

func WithErrorHandler(handler ErrorHandler) Option {
	return func(s *ServeMux) {
		s.errHandler = handler
	}
}

func WithPathPrefix(prefix string) Option {
	trimmedPrefix := strings.Trim(prefix, "/")

	return func(s *ServeMux) {
		s.pathFn = func(route string) string {
			return "/" + trimmedPrefix + "/" + route
		}
	}
}

func WithMiddleware(mw Middleware) Option {
	return func(s *ServeMux) {
		s.middleware = append(s.middleware, mw)
	}
}

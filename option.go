package ferry

import "net/http"

type Option func(mux *mux)

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

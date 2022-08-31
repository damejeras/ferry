package ferry

import "net/http"

func WithErrorHandler(handler ErrorHandler) func(*mux) {
	return func(m *mux) {
		m.errHandler = handler
	}
}

func WithNotFound(handler http.HandlerFunc) func(*mux) {
	return func(m *mux) {
		m.NotFound(handler)
	}
}

func WithMethodNotAllowed(handler http.HandlerFunc) func(*mux) {
	return func(m *mux) {
		m.MethodNotAllowed(handler)
	}
}

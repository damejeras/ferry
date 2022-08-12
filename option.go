package ferry

type Option func(mux *mux)

func WithErrorHandler(handler ErrorHandler) Option {
	return func(s *mux) {
		s.errHandler = handler
	}
}

func WithMiddleware(mw Middleware) Option {
	return func(s *mux) {
		s.middleware = append(s.middleware, mw)
	}
}

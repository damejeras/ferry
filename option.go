package ferry

type Option func(mux *mux)

func WithErrorHandler(handler ErrorHandler) Option {
	return func(s *mux) {
		s.errHandler = handler
	}
}

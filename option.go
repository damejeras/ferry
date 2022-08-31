package ferry

func WithErrorHandler(handler ErrorHandler) func(*mux) {
	return func(m *mux) {
		m.errHandler = handler
	}
}

package ferry

type Option func(mux *mux)

func WithErrorHandler(handler ErrorHandler) Option {
	return func(m *mux) {
		m.errHandler = handler
	}
}

func WithServiceDiscovery(m *mux) {
	m.serviceDiscovery = true
}

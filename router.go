package ferry

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

// Router is the wrapper of chi.Router which allows to register Procedure and Stream handlers.
// This router is intended to be mounted on regular chi Router.
type Router interface {
	// Register registers Procedure or Stream Handler.
	Register(...Handler)

	chi.Router
}

// NewRouter creates a Router instance. Router is the extension of chi.Router.
func NewRouter(options ...Option) Router {
	router := chi.NewRouter()

	m := &mux{
		errHandler: DefaultErrorHandler,
		handlers:   make(map[string]Handler),

		Router: router,
	}

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if err := Encode(w, r, http.StatusNotFound, ClientError{
			Code:    http.StatusNotFound,
			Message: "not found",
		}); err != nil {
			m.errHandler(w, r, err)
		}
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		if err := Encode(w, r, http.StatusNotFound, ClientError{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed",
		}); err != nil {
			m.errHandler(w, r, err)
		}
	})

	for i := range options {
		options[i](m)
	}

	return m
}

// mux is the implementation of Router interface.
type mux struct {
	errHandler ErrorHandler
	handlers   map[string]Handler
	mutex      sync.Mutex
	ready      bool

	serviceDiscovery bool

	chi.Router
}

// Register registers Procedure or Stream handlers to the Router.
func (m *mux) Register(handlers ...Handler) {
	for _, h := range handlers {
		if exists, ok := m.handlers[h.serviceMeta.methodName]; ok {
			panic(fmt.Sprintf(
				"can not register %q, because %q already exists; functions must not share %q",
				h.serviceMeta.reflectedName,
				exists.serviceMeta.reflectedName,
				h.serviceMeta.methodName,
			))
		}

		switch h.handlerType {
		case procedureHandler:
			m.Post(h.serviceMeta.path(), h.builder(m))
		case streamHandler:
			m.Get(h.serviceMeta.path(), h.builder(m))
		default:
			continue
		}

		m.handlers[h.serviceMeta.methodName] = h
	}
}

// ServeHTTP is the implementation of http.Handler interface.
func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.serviceDiscovery && !m.ready {
		m.init()
	}

	m.Router.ServeHTTP(w, r)
}

// init initializes and registers spec handler.
func (m *mux) init() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.ready {
		return
	}

	m.Router.Handle("/", specHandler(m.handlers))

	m.ready = true
}

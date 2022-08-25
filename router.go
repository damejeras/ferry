package ferry

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

// Router is the wrapper of chi.Router which allows to register Procedure and Stream handlers.
// This router is intended to be mounted on regular chi Router.
type Router interface {
	// Register registers Procedure or Stream Handler.
	Register(Handler)

	chi.Router
}

// NewRouter creates a Router instance. Router is the extension of chi.Router.
func NewRouter(options ...Option) Router {
	router := chi.NewRouter()

	m := &mux{
		errHandler: DefaultErrorHandler,
		spec:       make([]spec, 0),

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
	spec       []spec
	mutex      sync.Mutex
	ready      bool

	chi.Router
}

func (m *mux) Register(h Handler) {
	s, handle := h(m)

	switch s.handlerType {
	case procedureHandler:
		m.Post(s.path(), handle)
		m.spec = append(m.spec, s)
	case streamHandler:
		m.Get(s.path(), handle)
		m.spec = append(m.spec, s)
	default:
		return
	}
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.ready {
		m.init()
	}

	m.Router.ServeHTTP(w, r)
}

// init initializes spec handler on first request to the router.
func (m *mux) init() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.ready {
		return
	}

	m.Router.Get("/", specHandler(m.spec))
	m.ready = true
}

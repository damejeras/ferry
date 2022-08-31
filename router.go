package ferry

import (
	"net/http"

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
func NewRouter(options ...func(m *mux)) Router {
	router := chi.NewRouter()

	m := &mux{
		errHandler: DefaultErrorHandler,
		Router:     router,
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

	chi.Router
}

// Register registers Procedure or Stream handlers to the Router.
func (m *mux) Register(handlers ...Handler) {
	for _, handler := range handlers {
		handler.build(m)

		switch h := handler.(type) {
		case *procedureHandler:
			m.Method(http.MethodPost, "/"+h.meta.name, handler)
		case *streamHandler:
			m.Method(http.MethodGet, "/"+h.meta.name, handler)
		default:
			continue
		}
	}
}

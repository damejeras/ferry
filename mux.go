package ferry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Router is the wrapper of chi.Router which allows to register Procedure and Stream handlers.
// This router is intended to be mounted on regular chi Router.
type Router interface {
	Register(Handler)

	chi.Router
}

// Handler can only be acquired from helper methods (Procedure, Stream).
// It provides type safety when defining API.
type Handler interface {
	build(mux *mux) (string, http.HandlerFunc)
}

// NewRouter creates a Router instance. Router is the extension of chi.Router.
func NewRouter(options ...Option) Router {
	router := chi.NewRouter()

	m := &mux{
		errHandler: DefaultErrorHandler,

		Router: router,
	}

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if err := EncodeJSONResponse(w, r, http.StatusNotFound, ClientError{
			Code:    http.StatusNotFound,
			Message: "not found",
		}); err != nil {
			m.errHandler(w, r, err)
		}
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		if err := EncodeJSONResponse(w, r, http.StatusNotFound, ClientError{
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

type mux struct {
	errHandler ErrorHandler

	chi.Router
}

func (m *mux) Register(h Handler) {
	path, handle := h.build(m)

	switch h.(type) {
	case procedureBuilder:
		m.Post(path, handle)
	case streamBuilder:
		m.Get(path, handle)
	default:
		return
	}
}

package ferry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ServeMux is a replacement for http.ServeMux which allows registering only remote procedures and streams.
type Router interface {
	chi.Router

	Register(Handler)
}

// Handler can only be acquired from helper methods (Procedure, Stream).
// It provides type safety when defining API.
type Handler interface {
	build(mux *mux) (string, http.HandlerFunc)
}

func NewServeMux(options ...Option) Router {
	apiRouter := chi.NewRouter()

	mux := &mux{
		errHandler: DefaultErrorHandler,

		Router: apiRouter,
	}

	apiRouter.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := EncodeJSON(w, r, http.StatusNotFound, ClientError{
			Code:    http.StatusNotFound,
			Message: "not found",
		}); err != nil {
			mux.errHandler(w, r, err)
		}
	}))

	apiRouter.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := EncodeJSON(w, r, http.StatusNotFound, ClientError{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed",
		}); err != nil {
			mux.errHandler(w, r, err)
		}
	}))

	for i := range options {
		options[i](mux)
	}

	return mux
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

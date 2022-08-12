package ferry

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

type ServeMux interface {
	Handle(path string, h Handler)
	OpenAPISpec(modification func(*openapi3.Spec)) ([]byte, error)

	http.Handler
}

func NewServeMux(options ...Option) ServeMux {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithDescription("Put something here")

	mux := &mux{
		apiReflector: reflector,
		errHandler:   DefaultErrorHandler,
		middleware:   make([]Middleware, 0),
		procedures:   make(map[string]http.Handler),
		streams:      make(map[string]http.Handler),
	}

	mux.notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := EncodeJSON(w, r, http.StatusNotFound, ClientError{
			Code:    http.StatusNotFound,
			Message: "not found",
		}); err != nil {
			mux.errHandler(w, r, err)
		}
	})

	for i := range options {
		options[i](mux)
	}

	return mux
}

type mux struct {
	apiReflector    openapi3.Reflector
	errHandler      ErrorHandler
	middleware      []Middleware
	notFoundHandler http.Handler
	procedures      map[string]http.Handler
	streams         map[string]http.Handler
}

func (m *mux) Handle(path string, h Handler) {
	handle, err := h.build(path, m)
	if err != nil {
		// error could come from openapi reflector when reading request and response objects
		panic(err)
	}

	switch h.(type) {
	case procedureBuilder:
		m.procedures[path] = chainMiddleware(handle, m.middleware...)
	case streamBuilder:
		m.streams[path] = chainMiddleware(handle, m.middleware...)
	default:
		return
	}
}

func (m *mux) OpenAPISpec(modification func(spec *openapi3.Spec)) ([]byte, error) {
	if modification != nil {
		modification(m.apiReflector.Spec)
	}

	return m.apiReflector.Spec.MarshalJSON()
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	switch r.Method {
	case http.MethodPost:
		handler = m.procedures[r.URL.Path]
	case http.MethodGet:
		handler = m.streams[r.URL.Path]
	}

	if handler != nil {
		handler.ServeHTTP(w, r)
		return
	}

	m.notFoundHandler.ServeHTTP(w, r)
}

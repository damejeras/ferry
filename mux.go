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
	mux := &mux{
		apiReflector: newReflector(),
		errHandler:   DefaultErrorHandler,
		middleware:   make([]func(http.Handler) http.Handler, 0),
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
	middleware      []func(http.Handler) http.Handler
	notFoundHandler http.Handler
	procedures      map[string]http.Handler
	streams         map[string]http.Handler
}

func (m *mux) Handle(path string, h Handler) {
	url, err := parseURL(path)
	if err != nil {
		panic(err)
	}

	handle, err := h.build(url, m)
	if err != nil {
		panic(err)
	}

	switch h.(type) {
	case procedureBuilder:
		m.procedures[url.path] = chainMiddleware(handle, m.middleware...)
	case streamBuilder:
		m.streams[url.path] = chainMiddleware(handle, m.middleware...)
	default:
		return
	}
}

func (m *mux) OpenAPISpec(mod func(spec *openapi3.Spec)) ([]byte, error) {
	if mod != nil {
		mod(m.apiReflector.Spec)
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

func chainMiddleware(h http.Handler, mw ...func(http.Handler) http.Handler) http.Handler {
	for i := range mw {
		h = mw[len(mw)-1-i](h)
	}

	return h
}

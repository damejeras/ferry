package ferry

import (
	"net/http"
	"strings"
	"sync"

	"github.com/swaggest/openapi-go/openapi3"
)

// ServeMux is a replacement for http.ServeMux which allows registering only remote procedures and streams.
type ServeMux interface {
	Handle(h Handler)
	http.Handler
}

// Handler can only be acquired from helper methods (Procedure, Stream).
// It provides type safety when defining API.
type Handler interface {
	build(mux *mux) (string, http.Handler, error)
}

func NewServeMux(apiPrefix string, options ...Option) ServeMux {
	mux := &mux{
		apiPrefix:    "/" + strings.Trim(apiPrefix, "/"),
		apiReflector: newReflector(),
		errHandler:   DefaultErrorHandler,
		getHandlers:  make(map[string]http.Handler),
		middleware:   make([]func(http.Handler) http.Handler, 0),
		postHandlers: make(map[string]http.Handler),
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
	apiPrefix       string
	apiReflector    openapi3.Reflector
	errHandler      ErrorHandler
	middleware      []func(http.Handler) http.Handler
	notFoundHandler http.Handler
	openapiSpecMod  func(spec *openapi3.Spec)
	postHandlers    map[string]http.Handler
	getHandlers     map[string]http.Handler
	mutex           sync.Mutex
	ready           bool
}

func (m *mux) Handle(h Handler) {
	path, handle, err := h.build(m)
	if err != nil {
		panic(err)
	}

	switch h.(type) {
	case procedureBuilder:
		m.postHandlers[path] = handle
	case streamBuilder:
		m.getHandlers[path] = handle
	default:
		return
	}
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.ready {
		// before serving initialize OpenAPI handlers
		openapiHandlers(m)
	}

	var handler http.Handler
	switch r.Method {
	case http.MethodPost:
		handler = m.postHandlers[r.URL.Path]
	case http.MethodGet:
		handler = m.getHandlers[r.URL.Path]
	}

	if handler != nil {
		chainMiddleware(handler, m.middleware...).ServeHTTP(w, r)
		return
	}

	chainMiddleware(m.notFoundHandler, m.middleware...).ServeHTTP(w, r)
}

func chainMiddleware(h http.Handler, mw ...func(http.Handler) http.Handler) http.Handler {
	for i := range mw {
		h = mw[len(mw)-1-i](h)
	}

	return h
}

func buildPath(prefix, service, method string) string {
	return strings.Join([]string{prefix, strings.Join([]string{service, method}, ".")}, "/")
}

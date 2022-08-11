package ferry

import (
	"context"
	"net/http"
)

func NewServeMux(options ...Option) *ServeMux {
	mux := &ServeMux{
		errHandler: DefaultErrorHandler,
		middleware: make([]Middleware, 0),
		pathFn: func(route string) string {
			return "/" + route
		},
		routes: make(map[string]http.Handler),
	}

	mux.notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := Encode(w, r, http.StatusNotFound, ClientError{
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

type ServeMux struct {
	errHandler      ErrorHandler
	middleware      []Middleware
	notFoundHandler http.Handler
	pathFn          func(route string) string
	routes          map[string]http.Handler
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// RPC must be called with method POST.
		mux.notFoundHandler.ServeHTTP(w, r)

		return
	}

	handler, ok := mux.routes[r.URL.Path]
	if ok {
		handler.ServeHTTP(w, r)

		return
	}

	// No route found.
	mux.notFoundHandler.ServeHTTP(w, r)
}

func RegisterHandler[Request any, Response any](mux *ServeMux, route string, method func(ctx context.Context, r *Request) (*Response, error)) {
	mux.routes[mux.pathFn(route)] = chainMiddleware(httpHandler(mux, method), mux.middleware...)
}

func httpHandler[Request any, Response any](mux *ServeMux, serviceMethod func(ctx context.Context, r *Request) (*Response, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestValue Request
		if err := Decode(r, &requestValue); err != nil {
			mux.errHandler(w, r, err)
			return
		}

		response, err := serviceMethod(r.Context(), &requestValue)
		if err != nil {
			mux.errHandler(w, r, err)
			return
		}

		if err := Encode(w, r, http.StatusOK, response); err != nil {
			mux.errHandler(w, r, err)
			return
		}
	})
}

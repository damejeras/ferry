package ferry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewServeMux(options ...Option) *ServeMux {
	mux := &ServeMux{
		errHandler: DefaultErrorHandler,
		middleware: make([]Middleware, 0),
		pathFn: func(route string) string {
			return "/" + route
		},
		procedures: make(map[string]http.Handler),
		streams:    make(map[string]http.Handler),
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

type ServeMux struct {
	errHandler      ErrorHandler
	middleware      []Middleware
	notFoundHandler http.Handler
	pathFn          func(route string) string
	procedures      map[string]http.Handler
	streams         map[string]http.Handler
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	switch r.Method {
	case http.MethodPost:
		handler = mux.procedures[r.URL.Path]
	case http.MethodGet:
		handler = mux.streams[r.URL.Path]
	}

	if handler != nil {
		handler.ServeHTTP(w, r)
		return
	}

	mux.notFoundHandler.ServeHTTP(w, r)
}

func RegisterProcedure[Request any, Response any](mux *ServeMux, route string, procedure func(ctx context.Context, r *Request) (*Response, error)) {
	mux.procedures[mux.pathFn(route)] = chainMiddleware(procedureHandler(mux, procedure), mux.middleware...)
}

func RegisterStream[Request any, Response any](mux *ServeMux, route string, stream func(ctx context.Context, r *Request) (<-chan *Response, error)) {
	mux.streams[mux.pathFn(route)] = chainMiddleware(streamHandler(mux, stream), mux.middleware...)
}

func procedureHandler[Request any, Response any](mux *ServeMux, procedure func(ctx context.Context, r *Request) (*Response, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestValue Request
		if err := DecodeJSON(r, &requestValue); err != nil {
			mux.errHandler(w, r, err)
			return
		}

		response, err := procedure(r.Context(), &requestValue)
		if err != nil {
			mux.errHandler(w, r, err)
			return
		}

		if err := EncodeJSON(w, r, http.StatusOK, response); err != nil {
			mux.errHandler(w, r, err)
			return
		}
	})
}

func streamHandler[Request any, Response any](mux *ServeMux, stream func(ctx context.Context, r *Request) (<-chan *Response, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			mux.errHandler(w, r, ClientError{
				Code:    http.StatusBadRequest,
				Message: "connection does not support streaming",
			})
			return
		}

		var requestValue Request
		if err := DecodeQuery(r, &requestValue); err != nil {
			mux.errHandler(w, r, err)
			return
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ctx := r.Context()
		eventStream, err := stream(ctx, &requestValue)
		if err != nil {
			mux.errHandler(w, r, err)
			return
		}

		var buffer bytes.Buffer
		for event := range eventStream {
			if err := json.NewEncoder(&buffer).Encode(event); err != nil {
				mux.errHandler(w, r, fmt.Errorf("encode event: %w", err))
				return
			}

			if _, err := io.Copy(w, &buffer); err != nil {
				mux.errHandler(w, r, fmt.Errorf("copy buffer to response writer: %w", err))
				return
			}

			flusher.Flush()
		}
	})
}

package ferry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

// Handler can only be acquired from helper methods (Procedure, Stream).
// It provides type safety when defining API.
type Handler interface {
	build(path string, mux *mux) (http.Handler, error)
}

// Procedure will return Handler which can be used to register remote procedure in ServeMux
func Procedure[Request any, Response any](procedure func(ctx context.Context, r *Request) (*Response, error)) Handler {
	return procedureBuilder(func(path string, mux *mux) (http.Handler, error) {
		op := openapi3.Operation{}
		if err := mux.apiReflector.SetRequest(&op, new(Request), http.MethodPost); err != nil {
			return nil, err
		}
		if err := mux.apiReflector.SetJSONResponse(&op, new(Response), http.StatusOK); err != nil {
			return nil, err
		}
		if err := mux.apiReflector.SetJSONResponse(&op, new(ClientError), http.StatusBadRequest); err != nil {
			return nil, err
		}
		if err := mux.apiReflector.SetJSONResponse(&op, new(ServerError), http.StatusInternalServerError); err != nil {
			return nil, err
		}
		if err := mux.apiReflector.Spec.AddOperation(http.MethodPost, path, op); err != nil {
			return nil, err
		}

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
		}), nil
	})
}

// Stream will return Handler which can be used to register remote procedure in ServeMux
func Stream[Request any, Message any](stream func(ctx context.Context, r *Request) (<-chan *Message, error)) Handler {
	return streamBuilder(func(path string, mux *mux) (http.Handler, error) {
		op := openapi3.Operation{}
		if err := mux.apiReflector.SetRequest(&op, new(Request), http.MethodGet); err != nil {
			return nil, err
		}
		if err := mux.apiReflector.SetJSONResponse(&op, new(ClientError), http.StatusBadRequest); err != nil {
			return nil, err
		}
		if err := mux.apiReflector.SetJSONResponse(&op, new(ServerError), http.StatusInternalServerError); err != nil {
			return nil, err
		}
		if err := mux.apiReflector.Spec.AddOperation(http.MethodGet, path, op); err != nil {
			return nil, err
		}

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
		}), nil
	})
}

// procedureBuilder is the implementation of Handler interface.
type procedureBuilder func(path string, mux *mux) (http.Handler, error)

// streamBuilder is the implementation of handler interface
type streamBuilder func(path string, mux *mux) (http.Handler, error)

func (b procedureBuilder) build(path string, mux *mux) (http.Handler, error) { return b(path, mux) }
func (b streamBuilder) build(path string, mux *mux) (http.Handler, error)    { return b(path, mux) }

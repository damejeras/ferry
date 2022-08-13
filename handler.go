package ferry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Handler can only be acquired from helper methods (Procedure, Stream).
// It provides type safety when defining API.
type Handler interface {
	build(url *url, mux *mux) (http.Handler, error)
}

// Procedure will return Handler which can be used to register remote procedure in ServeMux
func Procedure[Req any, Res any](procedure func(ctx context.Context, r *Req) (*Res, error)) Handler {
	return procedureBuilder(func(url *url, mux *mux) (http.Handler, error) {
		op, err := procedureOp(url, mux, new(Req), new(Res))
		if err != nil {
			return nil, err
		}

		if err := mux.apiReflector.Spec.AddOperation(http.MethodPost, url.path, op); err != nil {
			return nil, err
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestValue Req
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
func Stream[Req any, Msg any](stream func(ctx context.Context, r *Req) (<-chan *Msg, error)) Handler {
	return streamBuilder(func(url *url, mux *mux) (http.Handler, error) {
		op, err := streamOp(url, mux, new(Req))
		if err != nil {
			return nil, err
		}

		if err := mux.apiReflector.Spec.AddOperation(http.MethodGet, url.path, op); err != nil {
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

			var reqValue Req
			if err := DecodeQuery(r, &reqValue); err != nil {
				mux.errHandler(w, r, err)
				return
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			ctx := r.Context()
			eventStream, err := stream(ctx, &reqValue)
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
type procedureBuilder func(url *url, mux *mux) (http.Handler, error)

// streamBuilder is the implementation of handler interface
type streamBuilder func(url *url, mux *mux) (http.Handler, error)

func (b procedureBuilder) build(url *url, mux *mux) (http.Handler, error) {
	return b(url, mux)
}

func (b streamBuilder) build(url *url, mux *mux) (http.Handler, error) {
	return b(url, mux)
}

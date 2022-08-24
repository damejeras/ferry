package ferry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// Event carries payload and event ID.
type Event[P any] struct {
	ID      string
	Payload *P
}

// Stream will return Handler which can be used to register SSE streams in ServeMux.
func Stream[Req any, Msg any](stream func(ctx context.Context, r *Req) (<-chan Event[Msg], error)) Handler {
	fn := runtime.FuncForPC(reflect.ValueOf(stream).Pointer()).Name()
	if !strings.HasSuffix(fn, "-fm") {
		panic("stream can only built from function with receiver")
	}

	nameParts := strings.Split(strings.TrimSuffix(fn, "-fm"), ".")
	if len(nameParts) < 2 {
		panic("stream can only built from function with receiver")
	}

	serviceName := nameParts[len(nameParts)-2]
	methodName := nameParts[len(nameParts)-1]

	return streamBuilder(func(mux *mux) (string, http.HandlerFunc) {
		path := "/" + serviceName + "." + methodName

		payloadType := reflect.TypeOf(new(Msg)).Elem().Name()

		handler := func(w http.ResponseWriter, r *http.Request) {
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
			events, err := stream(ctx, &reqValue)
			if err != nil {
				mux.errHandler(w, r, err)
				return
			}

			// respond immediately with keep-alive message
			if _, err := fmt.Fprintf(w, "event: keep-alive\n\n"); err != nil {
				mux.errHandler(w, r, fmt.Errorf("write initial keep-alive: %w", err))
				return
			}
			flusher.Flush()

			for {
				select {
				case event, ok := <-events:
					if !ok {
						return
					}
					payload, err := json.Marshal(event.Payload)
					if err != nil {
						mux.errHandler(w, r, fmt.Errorf("encode message: %w", err))
						return
					}
					if _, err := fmt.Fprintf(w, "id: %s\nevent: %s\ndata: %s\n\n", event.ID, payloadType, payload); err != nil {
						mux.errHandler(w, r, fmt.Errorf("write message: %w", err))
						return
					}
				// keep connection alive
				case <-time.After(5 * time.Second):
					if _, err := fmt.Fprintf(w, "event: keep-alive\n\n"); err != nil {
						mux.errHandler(w, r, fmt.Errorf("write keep-alive: %w", err))
						return
					}
				}

				flusher.Flush()
			}
		}

		return path, handler
	})
}

// streamBuilder is the implementation of handler interface
type streamBuilder func(mux *mux) (string, http.HandlerFunc)

func (b streamBuilder) build(mux *mux) (string, http.HandlerFunc) {
	return b(mux)
}

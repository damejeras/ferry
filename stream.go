package ferry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

// Event carries payload and event ID.
type Event[P any] struct {
	ID      string
	Payload *P
}

// Stream will return Handler which can be used to register SSE stream in Router.
// Stream function MUST close channel when context is cancelled. Handler will panic if context is cancelled and channel is not closed.
// Provided argument MUST be a function which has a receiver.
func Stream[Req any, Msg any](fn func(ctx context.Context, r *Req) (<-chan Event[Msg], error)) Handler {
	payloadType := reflect.TypeOf(new(Msg)).Elem().Name()

	meta, err := buildMeta(fn, new(Req))
	if err != nil {
		panic(err)
	}

	return Handler{
		handlerType: streamHandler,
		serviceMeta: meta,
		builder: func(m *mux) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				flusher, ok := w.(http.Flusher)
				if !ok {
					m.errHandler(w, r, ClientError{
						Code:    http.StatusBadRequest,
						Message: "connection does not support streaming",
					})
					return
				}

				var reqValue Req
				if err := decodeQuery(r, &reqValue); err != nil {
					m.errHandler(w, r, err)
					return
				}

				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				ctx := r.Context()
				events, err := fn(ctx, &reqValue)
				if err != nil {
					m.errHandler(w, r, err)
					return
				}

				// respond immediately with keep-alive message
				if _, err := fmt.Fprintf(w, "event: keep-alive\n\n"); err != nil {
					m.errHandler(w, r, fmt.Errorf("write initial keep-alive: %w", err))
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
							m.errHandler(w, r, fmt.Errorf("encode message: %w", err))
							return
						}
						if _, err := fmt.Fprintf(w, "id: %s\nevent: %s\ndata: %s\n\n", event.ID, payloadType, payload); err != nil {
							m.errHandler(w, r, fmt.Errorf("write message: %w", err))
							return
						}
					case <-time.After(5 * time.Second):
						select {
						case <-ctx.Done():
							// panic, channel MUST be closed when context is cancelled
							panic(fmt.Sprintf("%q stream channel is not closed", meta.methodName))
						default:
							// keep connection alive
							if _, err := fmt.Fprintf(w, "event: keep-alive\n\n"); err != nil {
								m.errHandler(w, r, fmt.Errorf("write keep-alive: %w", err))
								return
							}
						}
					}

					flusher.Flush()
				}
			}
		},
	}
}

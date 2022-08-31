package ferry

import (
	"net/http"
)

// Handler can only be acquired from helper methods (Procedure, Stream).
type Handler interface {
	http.Handler

	build(*mux)
}

type procedureHandler struct {
	builder func(m *mux) http.HandlerFunc
	http    func(http.ResponseWriter, *http.Request)

	meta
}

func (h *procedureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.http(w, r) }
func (h *procedureHandler) build(m *mux)                                     { h.http = h.builder(m) }

type streamHandler struct {
	builder     func(m *mux) http.HandlerFunc
	httpHandler func(http.ResponseWriter, *http.Request)

	meta
}

func (h *streamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.httpHandler(w, r) }
func (h *streamHandler) build(m *mux)                                     { h.httpHandler = h.builder(m) }

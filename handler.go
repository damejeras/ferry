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
	serviceMeta    serviceMeta
	handlerBuilder func(m *mux) http.HandlerFunc
	httpHandler    func(http.ResponseWriter, *http.Request)
}

func (h *procedureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.httpHandler(w, r) }
func (h *procedureHandler) build(m *mux)                                     { h.httpHandler = h.handlerBuilder(m) }

type streamHandler struct {
	serviceMeta    serviceMeta
	handlerBuilder func(m *mux) http.HandlerFunc
	httpHandler    func(http.ResponseWriter, *http.Request)
}

func (h *streamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.httpHandler(w, r) }
func (h *streamHandler) build(m *mux)                                     { h.httpHandler = h.handlerBuilder(m) }

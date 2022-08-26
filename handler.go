package ferry

import (
	"net/http"
	"strings"
)

type handlerType int

const (
	procedureHandler handlerType = iota
	streamHandler
)

// Handler can only be acquired from helper methods (Procedure, Stream).
// It provides type safety when defining API.
type Handler struct {
	handlerType handlerType
	serviceMeta serviceMeta
	builder     func(m *mux) http.HandlerFunc
}

// specHandler builds handler for the API.
func specHandler(handlers []Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointURL := "http://"
		if r.TLS != nil {
			endpointURL = "https://"
		}
		endpointURL += r.Host + r.URL.Path

		endpoints := make([]endpoint, 0)
		for i := range handlers {
			e := endpoint{
				Path:  strings.TrimSuffix(endpointURL, "/") + handlers[i].serviceMeta.path(),
				Body:  handlers[i].serviceMeta.body,
				Query: handlers[i].serviceMeta.query,
			}

			switch handlers[i].handlerType {
			case procedureHandler:
				e.Method = "POST"
			case streamHandler:
				e.Method = "GET"
			}

			endpoints = append(endpoints, e)
		}

		_ = RespondPretty(w, r, http.StatusOK, endpoints)
	}
}

type endpoint struct {
	Method string            `json:"method"`
	Path   string            `json:"path"`
	Body   map[string]string `json:"body,omitempty"`
	Query  map[string]string `json:"query,omitempty"`
}

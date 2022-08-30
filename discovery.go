package ferry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ServiceDiscovery(router chi.Router) http.HandlerFunc {
	endpoints := make([]endpoint, 0)

	chi.Walk(router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		switch h := handler.(type) {
		case *procedureHandler:
			endpoints = append(endpoints, endpoint{
				Method: method,
				Path:   route,
				Body:   h.serviceMeta.body,
				Query:  h.serviceMeta.query,
			})
		case *streamHandler:
			endpoints = append(endpoints, endpoint{
				Method: method,
				Path:   route,
				Body:   h.serviceMeta.body,
				Query:  h.serviceMeta.query,
			})
		}
		return nil
	})

	return func(w http.ResponseWriter, r *http.Request) {
		endpointURL := "http://"
		if r.TLS != nil {
			endpointURL = "https://"
		}
		endpointURL += r.Host

		_ = Respond(w, r, http.StatusOK, prependHost(endpointURL, endpoints))
	}
}

func prependHost(url string, input []endpoint) []endpoint {
	result := make([]endpoint, len(input))

	for i := range input {
		result[i] = endpoint{
			Method: input[i].Method,
			Path:   url + input[i].Path,
			Body:   input[i].Body,
			Query:  input[i].Query,
		}
	}

	return result
}

type endpoint struct {
	Method string            `json:"method"`
	Path   string            `json:"path"`
	Body   map[string]string `json:"body,omitempty"`
	Query  map[string]string `json:"query,omitempty"`
}

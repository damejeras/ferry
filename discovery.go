package ferry

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ServiceDiscovery walks chi.Router routing tree and creates http.HandlerFunc
// that will return list of ferry endpoints along with their metadata.
func ServiceDiscovery(router chi.Router) http.HandlerFunc {
	endpoints := make([]endpoint, 0)

	chi.Walk(router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		var m meta
		switch h := handler.(type) {
		case *procedureHandler:
			m = h.meta
		case *streamHandler:
			m = h.meta
		default:
			return nil
		}

		endpoints = append(endpoints, endpoint{
			Method: method,
			Path:   route,
			Body:   m.body,
			Query:  m.query,
		})

		return nil
	})

	return func(w http.ResponseWriter, r *http.Request) {
		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(prependHost(scheme+r.Host, endpoints))
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

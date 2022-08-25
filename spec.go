package ferry

import (
	"net/http"
	"strings"
)

type spec struct {
	httpMethod  string
	serviceName string
	methodName  string
	body        map[string]string
	query       map[string]string
}

func specHandler(spec []spec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}

		fullURL := scheme + r.Host + r.URL.Path

		endpoints := make([]endpoint, 0)

		for _, s := range spec {
			endpoints = append(endpoints, endpoint{
				Method: s.httpMethod,
				Path:   strings.TrimSuffix(fullURL, "/") + "/" + s.serviceName + "." + s.methodName,
				Body:   s.body,
				Query:  s.query,
			})
		}

		_ = IndentJSONResponse(w, r, http.StatusOK, endpoints)
	}
}

type endpoint struct {
	Method string            `json:"method"`
	Path   string            `json:"path"`
	Body   map[string]string `json:"body,omitempty"`
	Query  map[string]string `json:"query,omitempty"`
}

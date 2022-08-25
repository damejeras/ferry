package ferry

import (
	"net/http"
	"strings"
)

type spec struct {
	httpMethod  string
	serviceName string
	methodName  string
	params      map[string]string
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
			ns := endpoint{
				Method: s.httpMethod,
				Path:   strings.TrimSuffix(fullURL, "/") + "/" + s.serviceName + "." + s.methodName,
			}

			switch s.httpMethod {
			case http.MethodGet:
				ns.Query = s.params
			case http.MethodPost:
				ns.Body = s.params
			}

			endpoints = append(endpoints, ns)
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

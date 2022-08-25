package ferry

import (
	"net/http"
	"strings"
)

type spec struct {
	handlerType handlerType
	serviceName string
	methodName  string
	body        map[string]string
	query       map[string]string
}

func (s spec) path() string {
	return "/" + s.serviceName + "." + s.methodName
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
			e := endpoint{
				Path:  strings.TrimSuffix(fullURL, "/") + s.path(),
				Body:  s.body,
				Query: s.query,
			}

			switch s.handlerType {
			case procedureHandler:
				e.Method = http.MethodPost
			case streamHandler:
				e.Method = http.MethodGet
			}

			endpoints = append(endpoints, e)
		}

		_ = IndentEncode(w, r, http.StatusOK, endpoints)
	}
}

type endpoint struct {
	Method string            `json:"method"`
	Path   string            `json:"path"`
	Body   map[string]string `json:"body,omitempty"`
	Query  map[string]string `json:"query,omitempty"`
}

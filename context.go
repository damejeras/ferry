package ferry

import (
	"context"
	"net/http"
)

type contextKey int

const (
	// ResponseWriter references http.ResponseWriter
	ResponseWriter contextKey = iota
	// Request references *http.Request
	Request
)

func createContext(w http.ResponseWriter, r *http.Request) context.Context {
	return context.WithValue(context.WithValue(r.Context(), Request, r), ResponseWriter, w)
}

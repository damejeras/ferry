package ferry

import "net/http"

type handlerType int

const (
	procedureHandler handlerType = iota
	streamHandler
)

// Handler can only be acquired from helper methods (Procedure, Stream).
// It provides type safety when defining API.
type Handler func(mux *mux) (spec, http.HandlerFunc)

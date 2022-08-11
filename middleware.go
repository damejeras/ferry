package ferry

import "net/http"

type Middleware func(http.Handler) http.Handler

func chainMiddleware(handle http.Handler, mw ...Middleware) http.Handler {
	for i := range mw {
		handle = mw[len(mw)-1-i](handle)
	}

	return handle
}

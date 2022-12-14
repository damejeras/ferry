package ferry

import (
	"context"
	"net/http"
)

// Procedure will return Handler which can be used to register remote procedure in Router.
// This function call will panic if procedure function does not have receiver or Request structure is unparsable.
func Procedure[Req any, Res any](fn func(ctx context.Context, r *Req) (*Res, error)) Handler {
	mt, err := buildMeta(fn, new(Req))
	if err != nil {
		panic(err)
	}

	decodeFn := decodeJSON[Req]
	if len(mt.body) == 0 {
		// skip decoding if there are no parameters.
		decodeFn = func(r *http.Request, v *Req) error { return nil }
	}

	return &procedureHandler{
		meta: mt,
		builder: func(m *mux) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				var requestValue Req
				if err := decodeFn(r, &requestValue); err != nil {
					m.errHandler(w, r, err)
					return
				}

				response, err := fn(createContext(w, r), &requestValue)
				if err != nil {
					m.errHandler(w, r, err)
					return
				}

				if err := Encode(w, r, http.StatusOK, response); err != nil {
					m.errHandler(w, r, err)
					return
				}
			}
		},
	}
}

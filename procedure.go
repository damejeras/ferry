package ferry

import (
	"context"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Procedure will return Handler which can be used to register remote procedure in Router.
// This function call will panic if provided function does not have a receiver.
func Procedure[Req any, Res any](fn func(ctx context.Context, r *Req) (*Res, error)) Handler {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if !strings.HasSuffix(name, "-fm") {
		panic("procedure can only built from function with receiver")
	}

	nameParts := strings.Split(strings.TrimSuffix(name, "-fm"), ".")
	if len(nameParts) < 2 {
		panic("procedure can only built from function with receiver")
	}

	serviceName := nameParts[len(nameParts)-2]
	methodName := nameParts[len(nameParts)-1]

	decodeFn := decodeJSON[Req]
	if reflect.TypeOf(new(Req)).Elem().NumField() == 0 {
		// skip decoding if there are no parameters.
		decodeFn = func(r *http.Request, v *Req) error { return nil }
	}

	body, err := jsonMapping(new(Req))
	if err != nil {
		panic(err)
	}

	query, err := queryMapping(new(Req))
	if err != nil {
		panic(err)
	}

	s := spec{
		handlerType: procedureHandler,
		serviceName: serviceName,
		methodName:  methodName,
		body:        body,
		query:       query,
	}

	return func(mux *mux) (spec, http.HandlerFunc) {
		return s, func(w http.ResponseWriter, r *http.Request) {
			var requestValue Req
			if err := decodeFn(r, &requestValue); err != nil {
				mux.errHandler(w, r, err)
				return
			}

			response, err := fn(r.Context(), &requestValue)
			if err != nil {
				mux.errHandler(w, r, err)
				return
			}

			if err := Encode(w, r, http.StatusOK, response); err != nil {
				mux.errHandler(w, r, err)
				return
			}
		}
	}
}

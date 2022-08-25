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
func Procedure[Req any, Res any](procedure func(ctx context.Context, r *Req) (*Res, error)) Handler {
	fn := runtime.FuncForPC(reflect.ValueOf(procedure).Pointer()).Name()
	if !strings.HasSuffix(fn, "-fm") {
		panic("procedure can only built from function with receiver")
	}

	nameParts := strings.Split(strings.TrimSuffix(fn, "-fm"), ".")
	if len(nameParts) < 2 {
		panic("procedure can only built from function with receiver")
	}

	serviceName := nameParts[len(nameParts)-2]
	methodName := nameParts[len(nameParts)-1]

	decodeFn := DecodeJSON[Req]
	if reflect.TypeOf(new(Req)).Elem().NumField() == 0 {
		// skip decoding if there are no parameters.
		decodeFn = func(r *http.Request, v *Req) error { return nil }
	}

	params, err := requestJSONReflection(new(Req))
	if err != nil {
		panic(err)
	}

	return procedureBuilder(func(mux *mux) (spec, http.HandlerFunc) {
		return spec{httpMethod: "POST", serviceName: serviceName, methodName: methodName, params: params},
			func(w http.ResponseWriter, r *http.Request) {
				var requestValue Req
				if err := decodeFn(r, &requestValue); err != nil {
					mux.errHandler(w, r, err)
					return
				}

				response, err := procedure(r.Context(), &requestValue)
				if err != nil {
					mux.errHandler(w, r, err)
					return
				}

				if err := EncodeJSONResponse(w, r, http.StatusOK, response); err != nil {
					mux.errHandler(w, r, err)
					return
				}
			}
	})
}

// procedureBuilder is the implementation of Handler interface.
type procedureBuilder func(mux *mux) (spec, http.HandlerFunc)

func (b procedureBuilder) build(mux *mux) (spec, http.HandlerFunc) {
	return b(mux)
}

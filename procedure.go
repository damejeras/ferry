package ferry

import (
	"context"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Procedure will return Handler which can be used to register remote procedure in ServeMux
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

	decodeFn := DecodeJSON
	if reflect.TypeOf(new(Req)).Elem().NumField() == 0 {
		decodeFn = func(r *http.Request, v any) error { return nil }
	}

	return procedureBuilder(func(mux *mux) (string, http.HandlerFunc) {
		path := "/" + serviceName + "." + methodName

		handler := func(w http.ResponseWriter, r *http.Request) {
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

			if err := EncodeJSON(w, r, http.StatusOK, response); err != nil {
				mux.errHandler(w, r, err)
				return
			}
		}

		return path, handler
	})
}

// procedureBuilder is the implementation of Handler interface.
type procedureBuilder func(mux *mux) (string, http.HandlerFunc)

func (b procedureBuilder) build(mux *mux) (string, http.HandlerFunc) {
	return b(mux)
}

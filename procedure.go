package ferry

import (
	"context"
	"net/http"
)

// Procedure will return Handler which can be used to register remote procedure in ServeMux
func Procedure[Req any, Res any](procedure func(ctx context.Context, r *Req) (*Res, error)) Handler {
	return procedureBuilder(func(meta Meta, mux *mux) (http.Handler, error) {
		op, err := procedureOp(meta, mux, new(Req), new(Res))
		if err != nil {
			return nil, err
		}

		if err := mux.apiReflector.Spec.AddOperation(http.MethodPost, meta.Path, op); err != nil {
			return nil, err
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestValue Req
			if err := DecodeJSON(r, &requestValue); err != nil {
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
		}), nil
	})
}

// procedureBuilder is the implementation of Handler interface.
type procedureBuilder func(meta Meta, mux *mux) (http.Handler, error)

func (b procedureBuilder) build(meta Meta, mux *mux) (http.Handler, error) {
	return b(meta, mux)
}

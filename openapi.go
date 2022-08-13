package ferry

import (
	"net/http"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

func newReflector() openapi3.Reflector {
	reflector := openapi3.Reflector{}
	reflector.DefaultOptions = append(reflector.DefaultOptions, jsonschema.StripDefinitionNamePrefix("Ferry"))
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}

	return reflector
}

func procedureOp[Req any, Res any](path *url, mux *mux, req *Req, res *Res) (openapi3.Operation, error) {
	op := openapi3.Operation{Tags: []string{path.service}}

	if err := mux.apiReflector.SetRequest(&op, req, http.MethodPost); err != nil {
		return openapi3.Operation{}, err
	}
	if err := mux.apiReflector.SetJSONResponse(&op, res, http.StatusOK); err != nil {
		return openapi3.Operation{}, err
	}
	if err := mux.apiReflector.SetJSONResponse(&op, new(ClientError), http.StatusBadRequest); err != nil {
		return openapi3.Operation{}, err
	}
	if err := mux.apiReflector.SetJSONResponse(&op, new(ServerError), http.StatusInternalServerError); err != nil {
		return openapi3.Operation{}, err
	}

	return op, nil
}

func streamOp[Req any](path *url, mux *mux, res *Req) (openapi3.Operation, error) {
	op := openapi3.Operation{Tags: []string{path.service}}
	if err := mux.apiReflector.SetRequest(&op, res, http.MethodGet); err != nil {
		return openapi3.Operation{}, err
	}
	if err := mux.apiReflector.SetJSONResponse(&op, new(ClientError), http.StatusBadRequest); err != nil {
		return openapi3.Operation{}, err
	}
	if err := mux.apiReflector.SetJSONResponse(&op, new(ServerError), http.StatusInternalServerError); err != nil {
		return openapi3.Operation{}, err
	}

	return op, nil
}

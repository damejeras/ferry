package ferry

import (
	"net/http"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

func openapiHandlers(m *mux) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if !m.ready {
		if m.openapiSpecMod != nil {
			m.openapiSpecMod(m.apiReflector.Spec)
			yaml, err := m.apiReflector.Spec.MarshalYAML()
			if err != nil {
				panic(err)
			}

			json, err := m.apiReflector.Spec.MarshalJSON()
			if err != nil {
				panic(err)
			}

			m.getHandlers[buildPath(m.apiPrefix, "openapi", "yaml")] = http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/x-yaml")
					_, err := w.Write(yaml)
					if err != nil {
						m.errHandler(w, r, err)
					}
				},
			)

			m.getHandlers[buildPath(m.apiPrefix, "openapi", "json")] = http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					_, err := w.Write(json)
					if err != nil {
						m.errHandler(w, r, err)
					}
				},
			)
		}

		m.ready = true
	}
}

func newReflector() openapi3.Reflector {
	reflector := openapi3.Reflector{}
	reflector.DefaultOptions = append(reflector.DefaultOptions, jsonschema.StripDefinitionNamePrefix("Ferry"))
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}

	return reflector
}

func procedureOp[Req any, Res any](serviceName string, mux *mux, req *Req, res *Res) (openapi3.Operation, error) {
	op := openapi3.Operation{Tags: []string{serviceName}}

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

func streamOp[Req any](serviceName string, mux *mux, res *Req) (openapi3.Operation, error) {
	op := openapi3.Operation{Tags: []string{serviceName}}
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

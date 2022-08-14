module github.com/damejeras/ferry/example

go 1.18

require github.com/damejeras/ferry v0.0.0-20220812150200-855f29bea2a6

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/swaggest/openapi-go v0.2.20
)

require (
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/swaggest/jsonschema-go v0.3.36 // indirect
	github.com/swaggest/refl v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/damejeras/ferry => ../.

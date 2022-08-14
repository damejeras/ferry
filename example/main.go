package main

import (
	"log"
	"net/http"

	"github.com/damejeras/ferry"
	"github.com/damejeras/ferry/example/internal/greet"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/swaggest/openapi-go/openapi3"
)

func main() {
	mux := ferry.NewServeMux(
		// a prefix for the api
		"/api/v1",
		// this will enable /api/v1/openapi.json and /api/v1/openapi.yaml endpoints
		ferry.WithOpenApiSpec(func(spec *openapi3.Spec) {
			spec.Servers = []openapi3.Server{{URL: "http://localhost:7777"}}
			spec.Info.WithDescription("This is example on how to use ferry.")
			spec.Info.WithTitle("Example API")
		}),
		// permissive CORS
		ferry.WithMiddleware(cors.AllowAll().Handler),
		// log requests to console
		ferry.WithMiddleware(middleware.Logger),
		// use default logging functionality but log errors with standard logger
		ferry.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			switch err.(type) {
			case ferry.ClientError:
				ferry.DefaultErrorHandler(w, r, err)
			default:
				log.Printf("unexpeced error: %v", err)
				ferry.DefaultErrorHandler(w, r, err)
			}
		}),
	)

	// This returns GreetService implementation from api/v1/greet.go. GreetService is the description of our API.
	// The implementation is done in internal/greet/service.go
	greetSvc := greet.NewService()
	// POST http://localhost:7777/api/v1/GreetService.HelloWorld
	// the endpoint name is being reflected from the GreetService in api/v1/greet.go
	mux.Handle(ferry.Procedure(greetSvc.HelloWorld))
	// POST http://localhost:7777/api/v1/GreetService.HelloName
	// Content-Type: application/json
	// { "name": "Joe" }
	// the endpoint name is being reflected from the GreetService in api/v1/greet.go
	mux.Handle(ferry.Procedure(greetSvc.HelloName))
	// GET http://localhost:7777/api/v1/GreetService.StreamGreetings
	// This will start streaming SSE events
	// the endpoint name is being reflected from the GreetService in api/v1/greet.go
	mux.Handle(ferry.Stream(greetSvc.StreamGreetings))

	if err := http.ListenAndServe(":7777", mux); err != nil {
		log.Fatal(err)
	}
}

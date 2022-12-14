package main

import (
	"log"
	"net/http"

	"github.com/damejeras/ferry"
	"github.com/damejeras/ferry/example/internal/greet"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	v1 := ferry.NewRouter(
		// use default error handling but log errors with standard logger
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

	// All the modifications to Router should be made before registering handlers.
	v1.Use(
		// permissive CORS
		cors.AllowAll().Handler,
		// log requests to console
		middleware.Logger,
	)

	// This returns GreetService implementation from api/v1/greet.go. GreetService is the description of our API.
	// The implementation is done in internal/greet/service.go
	greetSvc := greet.NewService()
	v1.Register(
		// POST http://localhost:7777/api/v1/HelloWorld
		// the endpoint path is reflected from function name.
		ferry.Procedure(greetSvc.HelloWorld),
		// POST http://localhost:7777/api/v1/HelloName
		// Content-Type: application/json
		// { "name": "Joe" }
		// the endpoint path is reflected from function name.
		ferry.Procedure(greetSvc.HelloName),
		// GET http://localhost:7777/api/v1/StreamGreetings
		// This will start streaming SSE events
		// the endpoint path is reflected from function name.
		ferry.Stream(greetSvc.StreamGreetings),
	)

	router := chi.NewRouter()
	router.Mount("/api/v1/GreetService", v1)
	// Must be called last, because ferry.ServiceDiscovery is walking router paths.
	router.Handle("/api/v1", ferry.ServiceDiscovery(router))

	if err := http.ListenAndServe(":7777", router); err != nil {
		log.Fatal(err)
	}
}

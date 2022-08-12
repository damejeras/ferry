package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/damejeras/ferry"
	"github.com/damejeras/ferry/example/internal/greet"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggest/openapi-go/openapi3"
)

func main() {
	mux := ferry.NewServeMux(
		//ferry.WithMiddleware(auth.Middleware),
		ferry.WithMiddleware(middleware.Logger),
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

	greetSvc := greet.NewService()
	mux.Handle("/api/v1/Greetings.HelloWorld", ferry.Procedure(greetSvc.HelloWorld))
	mux.Handle("/api/v1/Greetings.HelloName", ferry.Procedure(greetSvc.HelloName))
	mux.Handle("/api/v1/Greetings.Stream", ferry.Stream(greetSvc.StreamGreetings))

	spec, err := mux.OpenAPISpec(func(spec *openapi3.Spec) {
		spec.Servers = []openapi3.Server{{URL: "http://localhost:7777"}}
		spec.Info.WithDescription("This is example on how to use ferry.")
		spec.Info.WithTitle("Example API")
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", spec)

	if err := http.ListenAndServe(":7777", mux); err != nil {
		log.Fatal(err)
	}
}

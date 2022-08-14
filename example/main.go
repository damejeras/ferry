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
		"/api/v1",
		ferry.WithOpenApiSpec(func(spec *openapi3.Spec) {
			spec.Servers = []openapi3.Server{{URL: "http://localhost:7777"}}
			spec.Info.WithDescription("This is example on how to use ferry.")
			spec.Info.WithTitle("Example API")
		}),
		ferry.WithMiddleware(cors.AllowAll().Handler),
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
	mux.Handle(ferry.Procedure(greetSvc.HelloWorld))
	mux.Handle(ferry.Procedure(greetSvc.HelloName))
	mux.Handle(ferry.Stream(greetSvc.StreamGreetings))

	if err := http.ListenAndServe(":7777", mux); err != nil {
		log.Fatal(err)
	}
}

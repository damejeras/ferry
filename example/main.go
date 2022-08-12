//go:generate swag init --dir ./,./internal/greet --parseDependency
package main

import (
	"log"
	"net/http"

	"github.com/damejeras/ferry"
	"github.com/damejeras/ferry/example/internal/greet"
)

// @title           ferry example API
// @version         1.0
// @description     This is example on how to use ferry.
// @license.name  	MIT
// @host      		localhost:7777
// @BasePath  		/api/v1
func main() {
	mux := ferry.NewServeMux(
		//ferry.WithMiddleware(auth.Middleware),
		ferry.WithPathPrefix("/api/v1"),
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

	ferry.RegisterProcedure(mux, "Greeting.HelloWorld", greetSvc.HelloWorld)
	ferry.RegisterProcedure(mux, "Greeting.HelloName", greetSvc.HelloName)
	ferry.RegisterStream(mux, "Greeting.Stream", greetSvc.StreamGreetings)

	if err := http.ListenAndServe(":7777", mux); err != nil {
		log.Fatal(err)
	}
}

package ferry

import (
	"context"
	"log"
	"net/http"
)

type helloWorldRequest struct{}

type helloWorldResponse struct {
	Message string
}

type helloNameRequest struct {
	Name string
}

type helloNameResponse struct {
	Message string
}

type greetingService struct{}

func (s *greetingService) HelloWorld(_ context.Context, _ *helloWorldRequest) (*helloWorldResponse, error) {
	return &helloWorldResponse{Message: "Hello World"}, nil
}

func (s *greetingService) HelloName(_ context.Context, request *helloNameRequest) (*helloNameResponse, error) {
	return &helloNameResponse{Message: "Hello " + request.Name}, nil
}

func main() {
	mux := NewServeMux(
		WithPathPrefix("/api/v1"),
		//WithMiddleware(func(next http.Handler) http.Handler {
		//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		authHeader := r.Header.Get("Authorization")
		//		if authHeader != "supersecret" {
		//			DefaultErrorHandler(w, r, ClientError{
		//				Code:    http.StatusUnauthorized,
		//				Message: "unauthorized",
		//			})
		//		}
		//
		//		next.ServeHTTP(w, r)
		//	})
		//}),
		WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			switch err.(type) {
			case ClientError:
				DefaultErrorHandler(w, r, err)
			default:
				log.Printf("unexpeced error: %v", err)
				DefaultErrorHandler(w, r, err)
			}
		}),
	)
	greetingService := new(greetingService)

	/**
	POST: http://localhost:7777/api/v1/Greeting.HelloWorld
	Content-Type: application/json
	{}
	*/
	RegisterHandler(mux, "Greeting.HelloWorld", greetingService.HelloWorld)
	/**
	POST: http://localhost:7777/api/v1/Greeting.HelloName
	Content-Type: application/json
	{"name": "Joe"}
	*/
	RegisterHandler(mux, "Greeting.HelloName", greetingService.HelloName)

	if err := http.ListenAndServe(":7777", mux); err != nil {
		log.Fatal(err)
	}
}

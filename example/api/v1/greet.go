package v1

import (
	"context"
)

type GreetService interface {
	HelloWorld(context.Context, *HelloWorldRequest) (*HelloWorldResponse, error)
	HelloName(context.Context, *HelloNameRequest) (*HelloNameResponse, error)
	StreamGreetings(context.Context, *StreamGreetingsRequest) (<-chan *Greeting, error)
}

type HelloWorldRequest struct{}

type HelloWorldResponse struct {
	Message string `json:"message"`
}

type HelloNameRequest struct {
	Name string `json:"name"`
}

type HelloNameResponse struct {
	Message string `json:"message"`
}

type StreamGreetingsRequest struct {
	Name string `query:"name"`
}

type Greeting struct {
	Message string
}

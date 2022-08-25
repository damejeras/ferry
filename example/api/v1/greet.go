package v1

import (
	"context"
	"time"

	"github.com/damejeras/ferry"
)

type GreetService interface {
	HelloWorld(context.Context, *HelloWorldRequest) (*HelloWorldResponse, error)
	HelloName(context.Context, *HelloNameRequest) (*HelloNameResponse, error)
	StreamGreetings(context.Context, *StreamGreetingsRequest) (<-chan ferry.Event[Greeting], error)
}

type HelloWorldRequest struct{}

type HelloWorldResponse struct {
	Message string `json:"message"`
}

type HelloNameRequest struct {
	Name string `json:"name,omitempty"`
}

type HelloNameResponse struct {
	Message string `json:"message"`
}

type StreamGreetingsRequest struct {
	Name string     `query:"name"`
	Time *time.Time `query:"name"`
}

type Greeting struct {
	Message string
}

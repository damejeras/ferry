package greet

import (
	"context"
	"strconv"
	"time"

	"github.com/damejeras/ferry"
	"github.com/damejeras/ferry/example/api/v1"
)

func NewService() v1.GreetService {
	return &service{}
}

type service struct{}

func (s *service) HelloWorld(_ context.Context, _ *v1.HelloWorldRequest) (*v1.HelloWorldResponse, error) {
	return &v1.HelloWorldResponse{Message: "Hello World"}, nil
}

func (s *service) HelloName(_ context.Context, request *v1.HelloNameRequest) (*v1.HelloNameResponse, error) {
	return &v1.HelloNameResponse{Message: "Hello " + request.Name}, nil
}

func (s *service) StreamGreetings(ctx context.Context, r *v1.StreamGreetingsRequest) (<-chan ferry.Event[v1.Greeting], error) {
	stream := make(chan ferry.Event[v1.Greeting])

	go func() {
		id := 0
		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ctx.Done():
				close(stream)
				return
			case <-ticker:
				id++
				stream <- ferry.Event[v1.Greeting]{
					ID:      strconv.Itoa(id),
					Payload: &v1.Greeting{Message: "Hello"},
				}
			}
		}
	}()

	return stream, nil
}

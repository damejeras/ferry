package greet

import (
	"context"
	"time"

	"ferry/example/api/v1"
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

func (s *service) StreamGreetings(ctx context.Context, r *v1.StreamGreetingsRequest) (<-chan *v1.Greeting, error) {
	stream := make(chan *v1.Greeting)

	go func() {
		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ctx.Done():
				close(stream)
				return
			case <-ticker:
				stream <- &v1.Greeting{Message: "hello " + r.Name}
			}
		}
	}()

	return stream, nil
}

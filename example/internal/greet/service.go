package greet

import (
	"context"
	"time"

	"github.com/damejeras/ferry/example/api/v1"
)

func NewService() v1.GreetService {
	return &service{}
}

type service struct{}

// HelloWorld says hello to the world.
// @Summary      Say "Hello World"
// @Description  say "Hello World"
// @Tags         Greeting
// @Accept       json
// @Produce      json
// @Success      200  {object}  v1.HelloWorldResponse
// @Failure      400  {object}  ferry.ClientError
// @Failure      500  {object}  ferry.ServerError
// @Router       /Greeting.HelloWorld [post]
func (s *service) HelloWorld(_ context.Context, _ *v1.HelloWorldRequest) (*v1.HelloWorldResponse, error) {
	return &v1.HelloWorldResponse{Message: "Hello World"}, nil
}

// HelloName says hello to you.
// @Summary      Say "Hello {Name}"
// @Description  say "Hello {Name}"
// @Tags         Greeting
// @Accept       json
// @Produce      json
// @Success      200  {object}  v1.HelloNameResponse
// @Failure      400  {object}  ferry.ClientError
// @Failure      500  {object}  ferry.ServerError
// @Router       /Greeting.HelloName [post]
func (s *service) HelloName(_ context.Context, request *v1.HelloNameRequest) (*v1.HelloNameResponse, error) {
	return &v1.HelloNameResponse{Message: "Hello " + request.Name}, nil
}

// StreamGreetings is very welcoming.
// @Summary      Sends greetings every second
// @Description  Opens SSE to send greetings.
// @Tags         Greeting
// @Accept       json
// @Produce      json
// @Failure      400  {object}  ferry.ClientError
// @Failure      500  {object}  ferry.ServerError
// @Router       /Greeting.StreamGreetings [get]
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

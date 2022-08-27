package greet

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/damejeras/ferry"
	v1 "github.com/damejeras/ferry/example/api/v1"
)

type service struct{}

func NewService() v1.GreetService {
	return &service{}
}

func (s *service) HelloWorld(_ context.Context, _ *v1.HelloWorldRequest) (*v1.HelloWorldResponse, error) {
	return &v1.HelloWorldResponse{Message: "Hello World"}, nil
}

func (s *service) HelloName(_ context.Context, r *v1.HelloNameRequest) (*v1.HelloNameResponse, error) {
	return &v1.HelloNameResponse{Message: fmt.Sprintf(randomFormat(), r.Name)}, nil
}

func (s *service) StreamGreetings(ctx context.Context, r *v1.StreamGreetingsRequest) (<-chan ferry.Event[v1.Greeting], error) {
	if r.Name == "" {
		return nil, ferry.ClientError{
			Code:    http.StatusBadRequest,
			Message: "name is required",
		}
	}

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
					Payload: &v1.Greeting{Message: fmt.Sprintf(randomFormat(), r.Name)},
				}
			}
		}
	}()

	return stream, nil
}

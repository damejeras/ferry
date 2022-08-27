# ferry

### What is `ferry`?
`ferry` is minimal RPC framework for quick HTTP API prototyping. It uses generics and reflection to reduce boilerplate.
`ferry` wraps [chi router](https://github.com/go-chi/chi), so it is easy to mount it to applications that are using go-chi.
It handles RPC calls through HTTP POST requests and provides SSE streams compatable with javascript EventSource class.

### What `ferry` is not?
`ferry` is not full pledged framework and does not aim to be one. It abstracts away boilerplate, but it might be not suitable for very complex APIs.

### How to use `ferry`?
With ferry you can prototype APIs with Go code. Think `gRPC` but without `protobuf`. You can see full example [here](https://github.com/damejeras/ferry/tree/main/_example).

First you have to define you API spec:

```go
package v1

import (...)

type GreetService interface {
	HelloWorld(context.Context, *HelloWorldRequest) (*HelloWorldResponse, error)
}

type HelloWorldRequest struct{}

type HelloWorldResponse struct {
	Message string `json:"message"`
}
```
Now you have to implement `GreetService` interface.
```go
package greet

import (...)

type service struct{}

func NewService() v1.GreetService { return &service{} }

func (s *service) HelloWorld(_ context.Context, _ *v1.HelloWorldRequest) (*v1.HelloWorldResponse, error) {
	return &v1.HelloWorldResponse{Message: "Hello World"}, nil
}
```
Now you have to create `ferry.Router` and register your `HelloWorld` procedure:
```go
v1 := ferry.NewRouter()
greetSvc := greet.NewService()
v1.Register(Procedure(greetSvc.HelloWorld))
```
`v1` here is also `chi.Router`, so you can add middleware or mount it to some other router.
`ferry` is using reflection and generics on your API spec to create routes. It is also providing simple API discovery through root path of the router.
In this case that would be `GET` to `/`, if you visit it you get the response
```json
[
  {
    "method": "POST",
    "path": "http://example.com/GreetService.HelloWorld",
  }
]
```




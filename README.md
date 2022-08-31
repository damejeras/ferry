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
// create ferry service router
v1greet:= ferry.NewRouter()
// create instance of your service
greetSvc := greet.NewService()
// register your service method
v1greet.Register(ferry.Procedure(greetSvc.HelloWorld))

// create root router
chiRouter := chi.NewRouter()
// mount your ferry service router
chiRouter.Mount("/api/v1/GreetService", v1greet)
// enable service discovery (optional)
chiRouter.Handle("/api/v1", ferry.ServiceDiscovery(chiRouter))
// run your server
http.ListenAndServe(":7777", chiRouter)
```

That's it. Because `ferry.Router` has `chi.Router` embedded you can use all the nice things `chi` provides.

### Service Discovery

`ferry`'s service discovery is meant to be read by humans first. Handler for service discovery is created by walking
router's routing tree. If you enabled it, `/api/v1` response should like this:
```json
[
  {
    "method": "POST",
    "path": "http://localhost:7777/api/v1/GreetService/HelloWorld"
  }
]
```

Service discovery can also print request parameters if your request has properties with `query` or `json` tags.
Try changing `HelloWorldRequest` in your spec to:
```go
type HelloWorldRequest struct{
  Name string `json:"name"`
}
```
Now service discovery should look like this:
```json
[
  {
    "method": "POST",
    "path": "http://localhost:7777/api/v1/GreetService/HelloWorld",
    "body": {
      "name": "string"
    }
  }
]
```

### Server-Sent Events

`ferry` also supports SSE streams. To learn more, check out [example application](https://github.com/damejeras/ferry/tree/main/_example).

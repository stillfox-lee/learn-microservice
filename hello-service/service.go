package main

import "context"

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

type HelloService interface {
	Hello(ctx context.Context, name string) (HelloResponse, error)
}

type helloService struct{}

func (s *helloService) Hello(ctx context.Context, name string) (HelloResponse, error) {
	return HelloResponse{Message: "Hello, " + name}, nil
}

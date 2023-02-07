package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeHelloEndpoint(svc HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(HelloRequest)
		res, err := svc.Hello(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

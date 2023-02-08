package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
)

func decodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return HelloRequest{Name: r.FormValue("name")}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type limitMiddleware struct {
	timer time.Duration
	burst int
}

func (l limitMiddleware) wrap(e endpoint.Endpoint) endpoint.Endpoint {
	e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(l.timer), l.burst))(e)
	return e
}

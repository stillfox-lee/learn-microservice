package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func decodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return HelloRequest{Name: r.FormValue("name")}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

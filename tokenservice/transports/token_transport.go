package transports

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/stillfox-lee/learn-microservice/token/endpoints"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

func DecodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tokenID := r.FormValue("id")

	if tokenID == "" {
		return nil, ErrorBadRequest
	}
	return &endpoints.GetTokenRequest{
		ID: tokenID,
	}, nil
}

func DecodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func DecodeHealthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &endpoints.HealthRequest{}, nil
}

func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

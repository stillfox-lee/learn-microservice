package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/stillfox-lee/learn-microservice/token/services"
)

type TokenEndpoints struct {
	CreateTokenEndpoint endpoint.Endpoint
	GetTokenEndpoint    endpoint.Endpoint
	HealthEndpoint      endpoint.Endpoint
}

type HealthRequest struct{}

type HealthResponse struct {
	Status string `json:"status"`
}

func MakeHealthEndpoint(svc services.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.Health(ctx)
		return HealthResponse{Status: status}, nil
	}
}

type TokenDTO struct {
	ID         string
	ExpireTime int64
	Data       interface{}
}

type CreateTokenRequest struct {
	Data interface{} `json:"data"`
	TTL  int64       `json:"ttl"`
}

type CreateTokenResponse struct {
	ID string `json:"id"`
}

func MakeCreateTokenEndpoint(svc services.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		dto := request.(*CreateTokenRequest)
		tokenID, err := svc.CreateToken(ctx, dto.TTL, dto.Data)
		// FIXME: biz error should wrapper into response
		return &CreateTokenResponse{ID: tokenID}, err
	}
}

type GetTokenRequest struct {
	ID string `json:"id"`
}

type GetTokenResponse struct {
	Token *TokenDTO `json:"token"`
	Err   error     `json:"err"`
}

func MakeGetTokenEndpoint(svc services.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*GetTokenRequest)
		token, err := svc.GetToken(ctx, req.ID)
		if err != nil {
			return &GetTokenResponse{
				Token: nil,
				Err:   err,
			}, nil
		}
		dto := TokenDTO{
			ID:         token.ID,
			ExpireTime: token.ExpireTimestamp,
			Data:       token.Data,
		}
		return &GetTokenResponse{Token: &dto, Err: err}, nil
	}
}

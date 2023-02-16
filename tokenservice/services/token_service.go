package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-kit/log"

	"github.com/stillfox-lee/learn-microservice/token/model"
)

type TokenService interface {
	CreateToken(ctx context.Context, ttl int64, data interface{}) (string, error)
	GetToken(ctx context.Context, tokenID string) (*model.TokenEntity, error)
	Health(ctx context.Context) string
}

type TokenSvc struct {
	TokenDao *model.TokenDao
	Log      log.Logger
}

func (s *TokenSvc) CreateToken(ctx context.Context, ttl int64, data interface{}) (string, error) {
	// create unique tokenID
	var tokenID string
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("create token id fail: %s", err)
	}
	tokenID = base64.URLEncoding.EncodeToString(b)
	token := &model.TokenEntity{
		ID:              tokenID,
		ExpireTimestamp: time.Now().Unix() + ttl,
		Data:            data,
	}
	err := s.TokenDao.CreateToken(ctx, token)
	if err != nil {
		return "", err
	}
	return tokenID, nil
}

func (s *TokenSvc) GetToken(ctx context.Context, tokenID string) (*model.TokenEntity, error) {
	token, err := s.TokenDao.GetToken(ctx, tokenID)
	if err != nil {
		return nil, err

	}
	return token, nil
}

func (s *TokenSvc) Health(ctx context.Context) string {
	return "ok"
}

func MakeTokenService(tokenDao *model.TokenDao, log log.Logger) TokenService {
	return &TokenSvc{
		Log:      log,
		TokenDao: tokenDao,
	}
}

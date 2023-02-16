package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type TokenEntity struct {
	ID              string
	ExpireTimestamp int64
	Data            interface{} // json
}

// DataToString format data into string
func (t *TokenEntity) DataToString() string {
	d, err := json.Marshal(t.Data)
	if err != nil {
		return ""
	}
	return string(d)
}

func (t *TokenEntity) DataFromString(data string) error {
	return json.Unmarshal([]byte(data), &t.Data)
}

func (t *TokenEntity) ExpireDuration() time.Duration {
	expireTime := time.Unix(t.ExpireTimestamp, 0)
	return time.Until(expireTime)
}

func (t *TokenEntity) UpdateExpire(ttl int64) {
	t.ExpireTimestamp += ttl
}

func (t *TokenEntity) IsExpired() bool {
	expireTime := time.Unix(t.ExpireTimestamp, 0)
	return expireTime.Before(time.Now())
}

type TokenDao struct {
	redis *redis.Client
}

func MakeDao(r *redis.Client) *TokenDao {
	return &TokenDao{redis: r}
}

// CreateToken token's ID should not exist
func (d *TokenDao) CreateToken(ctx context.Context, entity *TokenEntity) error {
	bcmd := d.redis.SetNX(ctx, entity.ID, entity.DataToString(), entity.ExpireDuration())
	return bcmd.Err()
}

func (d *TokenDao) GetToken(ctx context.Context, tokenID string) (*TokenEntity, error) {
	ttl, err := d.redis.TTL(ctx, tokenID).Result()
	switch {
	case err == redis.Nil:
		return nil, fmt.Errorf("token %s not exist", tokenID)
	case err != nil:
		return nil, fmt.Errorf("get token fail: %s", err)
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("token %s not exist", tokenID)
	}
	resData, err := d.redis.Get(ctx, tokenID).Result()
	var data interface{}
	if resData != "" {
		err = json.Unmarshal([]byte(resData), &data)
		if err != nil {
			return nil, err
		}
	}
	token := &TokenEntity{
		ID:              tokenID,
		Data:            data,
		ExpireTimestamp: time.Now().Unix() + int64(ttl.Seconds()),
	}

	return token, err
}

// UpdateToken should update an exixt token who identity by ID
func (d *TokenDao) UpdateToken(ctx context.Context, token *TokenEntity) error {
	return nil
}

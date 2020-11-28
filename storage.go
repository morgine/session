package session

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Storage interface {
	// save token
	SaveToken(id, token string, expires int64) error

	// check and refresh token if token is OK
	CheckAndRefreshToken(id, token string, expires int64) (ok bool, err error)

	// remove token
	DelToken(id string) error
}

type RedisStorage struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisStorage(keyPrefix string, client *redis.Client) Storage {
	return &RedisStorage{keyPrefix: keyPrefix, client: client}
}

func (rs *RedisStorage) key(id string) string {
	return rs.keyPrefix + id
}

func (rs *RedisStorage) SaveToken(id, token string, expires int64) error {
	return rs.client.Set(context.Background(), rs.key(id), token, time.Duration(expires)*time.Second).Err()
}

func (rs *RedisStorage) CheckAndRefreshToken(id, token string, expires int64) (ok bool, err error) {
	key := rs.key(id)
	savedToken, err := rs.client.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	if savedToken != token {
		return false, nil
	} else {
		err = rs.client.Expire(context.Background(), key, time.Duration(expires)*time.Second).Err()
		if err != nil {
			return false, nil
		} else {
			return true, nil
		}
	}
}

func (rs *RedisStorage) DelToken(id string) error {
	return rs.client.Del(context.Background(), rs.key(id)).Err()
}

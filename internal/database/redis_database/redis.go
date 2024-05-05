package redis_database

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisRepositoryDb struct {
	redisClient *redis.Client
}

func NewRedisRepositoryDb(ctx context.Context, redis *redis.Client) *RedisRepositoryDb {
	return &RedisRepositoryDb{redisClient: redis}
}

func (rd *RedisRepositoryDb) MakeKey(ip, token string) string {
	return ip + ":" + token
}

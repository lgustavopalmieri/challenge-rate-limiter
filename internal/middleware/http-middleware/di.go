package middleware

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/database/redis_database"
)

func NewWebRateLimiterMiddleware(redisClient *redis.Client) *IpTokenRateLimiterMiddleware{
	rateLimiterRepository := redis_database.NewRedisRepositoryDb(context.Background(),redisClient)
	return NewWebIpTokenMiddleware(rateLimiterRepository)
}

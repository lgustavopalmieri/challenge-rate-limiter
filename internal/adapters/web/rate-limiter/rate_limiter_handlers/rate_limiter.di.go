package rate_limiter_handlers

import (
	"context"
	"ratelimiter2/internal/database/redis_database"

	"github.com/go-redis/redis/v8"
)

func NewWebRateLimiterMiddleware(redisClient *redis.Client) *RateLimiterHandler {
	rateLimiterRepository := redis_database.NewRedisRepositoryDb(context.Background(), redisClient)
	return NewWebRateLimiterHandlers(rateLimiterRepository)
}

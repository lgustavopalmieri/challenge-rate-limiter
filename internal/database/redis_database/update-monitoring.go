package redis_database

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_entity"
)

func (rd *RedisRepositoryDb) UpdateMonitoring(ctx context.Context, rateLimiter rate_limiter_entity.RateLimiter) error {
	rateLimiterJSON, err := json.Marshal(rateLimiter)
	if err != nil {
		log.Println("failed to convert RateLimiter to JSON")
		return err
	}
	key := rd.MakeKey(rateLimiter.IP, rateLimiter.Token)

	if err := rd.redisClient.Set(ctx, key, rateLimiterJSON, 0).Err(); err != nil {
		log.Println("failed to update data in Redis")
		return err
	}

	return err
}

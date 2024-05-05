package redis_database

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_entity"
)

func (rd *RedisRepositoryDb) FindLimiter(ctx context.Context, ip, token string) (*rate_limiter_entity.RateLimiter, error) {
	log.Println("bateu find redis", ip, token)

	key := rd.MakeKey(ip, token)
	rateLimiterJSON, err := rd.redisClient.Get(ctx, key).Result()
	log.Println(rateLimiterJSON)

	if err != nil {
		return nil, err
	}
	var rateLimiter rate_limiter_entity.RateLimiter
	if err := json.Unmarshal([]byte(rateLimiterJSON), &rateLimiter); err != nil {
		log.Println("failed to parse rate limiter to JSON")
		return nil, err
	}
	return &rateLimiter, nil
}

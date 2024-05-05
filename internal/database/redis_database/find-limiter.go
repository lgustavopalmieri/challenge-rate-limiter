package redis_database

import (
	"context"
	"encoding/json"
	"log"
	"ratelimiter2/internal/rate-limiter/rate_limiter_entity"
)

func (rd *RedisRepositoryDb) FindLimiter(ctx context.Context, ip, token string) (*rate_limiter_entity.RateLimiter, error) {
	key := rd.MakeKey(ip, token)

	rateLimiterJSON, err := rd.redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Erro ao buscar rate limiter do Redis: %v\n", err)
		return nil, err
	}

	var rateLimiter rate_limiter_entity.RateLimiter
	if err := json.Unmarshal([]byte(rateLimiterJSON), &rateLimiter); err != nil {
		log.Println("Falha ao analisar rate limiter do JSON:", err)
		return nil, err
	}

	return &rateLimiter, nil
}

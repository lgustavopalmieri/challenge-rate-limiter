package redis_database

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_entity"
)

func (rd *RedisRepositoryDb) InitMonitoring(ctx context.Context, rateLimiter *rate_limiter_entity.RateLimiter) error {

	if rateLimiter == nil {
		return errors.New("rateLimiter is null")
	}

	// rateLimiterEntity := &rate_limiter_entity.RateLimiter{
	// 	IP:                     rateLimiter.IP,
	// 	Token:                  rateLimiter.Token,
	// 	IPLimit:                rateLimiter.IPLimit,
	// 	TokenLimit:             rateLimiter.TokenLimit,
	// 	BlockDurationInSeconds: rateLimiter.BlockDurationInSeconds,
	// 	InitTryingAt:           rateLimiter.InitTryingAt,
	// 	LastTryingAt:           rateLimiter.LastTryingAt,
	// 	Authorized:             rateLimiter.Authorized,
	// }

	jsonData, err := json.Marshal(rateLimiter)
	if err != nil {
		log.Println("failed to serialize rateLimiterEntity to JSON:", err)
		return err
	}

	key := rd.MakeKey(rateLimiter.IP, rateLimiter.Token)

	err = rd.redisClient.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		log.Println("fail saving on redis")
		return err
	}

	return err
}

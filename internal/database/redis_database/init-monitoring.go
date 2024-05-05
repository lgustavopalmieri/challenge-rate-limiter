package redis_database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"ratelimiter2/internal/rate-limiter/rate_limiter_entity"
	"time"
)

func (rd *RedisRepositoryDb) InitMonitoring(ctx context.Context, rateLimiter *rate_limiter_entity.RateLimiter) error {
	if rateLimiter == nil {
		return errors.New("rateLimiter is nil")
	}

	expirationTimeStr := os.Getenv("EXPIRATION_TIME")
	if expirationTimeStr == "" {
		return errors.New("EXPIRATION_TIME not set in environment variables")
	}

	expirationTime, err := time.ParseDuration(expirationTimeStr + "s")
	if err != nil {
		return fmt.Errorf("failed to parse EXPIRATION_TIME: %v", err)
	}

	jsonData, err := json.Marshal(rateLimiter)
	if err != nil {
		log.Println("Failed to serialize rateLimiter to JSON:", err)
		return err
	}

	key := rd.MakeKey(rateLimiter.IP, rateLimiter.Token)

	if err := rd.redisClient.Set(ctx, key, jsonData, expirationTime * 2).Err(); err != nil {
		log.Printf("Failed to save data to Redis for key %s: %v\n", key, err)
		return err
	}

	return nil
}

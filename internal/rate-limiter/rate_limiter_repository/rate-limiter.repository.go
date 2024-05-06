package rate_limiter_repository

import (
	"context"
	"ratelimiter2/internal/rate-limiter/rate_limiter_entity"
)

type RateLimiterRepositoryInterface interface {
	InitMonitoring(ctx context.Context, rateLimiter *rate_limiter_entity.RateLimiter) error
	FindLimiter(ctx context.Context, ip, token string) (*rate_limiter_entity.RateLimiter, error)
}

package rate_limiter_repository

import (
	"context"

	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_entity"
)

type RateLimiterRepositoryInterface interface {
	InitMonitoring(ctx context.Context, rateLimiter *rate_limiter_entity.RateLimiter) (*rate_limiter_entity.RateLimiter, error)
	FindLimiter(ctx context.Context, ip, token string) (*rate_limiter_entity.RateLimiter, error)
	UpdateMonitoring(ctx context.Context, rateLimiter rate_limiter_entity.RateLimiter) (*rate_limiter_entity.RateLimiter, error)
	CleanLimiter(ctx context.Context, ip, token string) error
}

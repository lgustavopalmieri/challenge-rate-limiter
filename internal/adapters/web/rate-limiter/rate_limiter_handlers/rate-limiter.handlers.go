package rate_limiter_handlers

import "ratelimiter2/internal/rate-limiter/rate_limiter_repository"

type RateLimiterHandler struct {
	RateLimiterRepository rate_limiter_repository.RateLimiterRepositoryInterface
}

func NewWebRateLimiterHandlers(
	rateLimiterRepository rate_limiter_repository.RateLimiterRepositoryInterface,

) *RateLimiterHandler {
	return &RateLimiterHandler{
		RateLimiterRepository: rateLimiterRepository,
	}
}

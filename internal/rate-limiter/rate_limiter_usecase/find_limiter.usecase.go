package rate_limiter_usecase

import (
	"context"
	"ratelimiter2/internal/rate-limiter/rate_limiter_entity"
	"ratelimiter2/internal/rate-limiter/rate_limiter_repository"
	"time"
)

type FindLimiterUseCase struct {
	rateLimiterRepository rate_limiter_repository.RateLimiterRepositoryInterface
}

func NewFindLimiterUseCase(ctx context.Context, rateLimiterRepository rate_limiter_repository.RateLimiterRepositoryInterface) *FindLimiterUseCase {
	return &FindLimiterUseCase{
		rateLimiterRepository: rateLimiterRepository,
	}
}

type FindLimiterOutputDTO struct {
	IP                     string    `json:"ip"`
	Token                  string    `json:"token,omitempty"`
	IPLimit                int64     `json:"ip_limit"`
	TokenLimit             int64     `json:"token_limit"`
	BlockDurationInSeconds int64     `json:"block_duration_in_seconds"`
	Reqs                   int64     `json:"reqs,omitempty"`
	InitTryingAt           time.Time `json:"init_trying_at"`
	LastTryingAt           time.Time `json:"last_trying_at"`
	Authorized             bool      `json:"authorized"`
}

func (uc *FindLimiterUseCase) Execute(ctx context.Context, ip, token string) (*FindLimiterOutputDTO, error) {
	newIp := rate_limiter_entity.RemoveIpPort(ip)
	existsLimiter, err := uc.rateLimiterRepository.FindLimiter(ctx, newIp, token)
	if err != nil {
		return nil, err
	}
	return &FindLimiterOutputDTO{
		IP:                     existsLimiter.IP,
		Token:                  existsLimiter.Token,
		IPLimit:                existsLimiter.IPLimit,
		TokenLimit:             existsLimiter.TokenLimit,
		BlockDurationInSeconds: existsLimiter.BlockDurationInSeconds,
		Reqs:                   existsLimiter.Reqs,
		InitTryingAt:           existsLimiter.InitTryingAt,
		LastTryingAt:           existsLimiter.LastTryingAt,
		Authorized:             existsLimiter.Authorized,
	}, nil
}

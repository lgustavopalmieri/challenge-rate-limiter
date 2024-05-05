package rate_limiter_usecase

import (
	"context"
	"log"
	"ratelimiter2/internal/rate-limiter/rate_limiter_entity"
	"ratelimiter2/internal/rate-limiter/rate_limiter_repository"
	"sync"
	"time"
)

type InitMonitoringUseCase struct {
	rateLimiterRepository rate_limiter_repository.RateLimiterRepositoryInterface
}

func NewInitMonitoringUseCase(ctx context.Context, rateLimiterRepository rate_limiter_repository.RateLimiterRepositoryInterface) *InitMonitoringUseCase {
	return &InitMonitoringUseCase{
		rateLimiterRepository: rateLimiterRepository,
	}
}

type InitMonitoringInputDTO struct {
	IP    string `json:"ip"`
	Token string `json:"token,omitempty"`
}

type InitMonitoringOutputDTO struct {
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

func (uc *InitMonitoringUseCase) Execute(ctx context.Context, input InitMonitoringInputDTO) (*InitMonitoringOutputDTO, error) {
	var mutex sync.Mutex
	newIp := rate_limiter_entity.RemoveIpPort(input.IP)

	mutex.Lock()
	existingLimiter, err := uc.rateLimiterRepository.FindLimiter(ctx, newIp, input.Token)
	if err != nil {
		log.Println("Erro ao buscar rate limiter do Redis:", err)
	}
	mutex.Unlock()
	if existingLimiter != nil {
		mutex.Lock()
		up := rate_limiter_entity.RateLimiter{
			IP:                     existingLimiter.IP,
			Token:                  existingLimiter.Token,
			IPLimit:                existingLimiter.IPLimit,
			TokenLimit:             existingLimiter.TokenLimit,
			BlockDurationInSeconds: existingLimiter.BlockDurationInSeconds,
			Reqs:                   existingLimiter.Reqs,
			InitTryingAt:           existingLimiter.InitTryingAt,
			LastTryingAt:           existingLimiter.LastTryingAt,
			Authorized:             existingLimiter.Authorized,
		}
		updated := rate_limiter_entity.UpdateLimiter(up)
		err = uc.rateLimiterRepository.InitMonitoring(ctx, updated)
		if err != nil {
			log.Println("Erro ao atualizar o rate limiter no Redis:", err)
		}
		mutex.Unlock()
		return &InitMonitoringOutputDTO{
			IP:                     updated.IP,
			Token:                  updated.Token,
			IPLimit:                updated.IPLimit,
			TokenLimit:             updated.TokenLimit,
			BlockDurationInSeconds: updated.BlockDurationInSeconds,
			Reqs:                   updated.Reqs,
			InitTryingAt:           updated.InitTryingAt,
			LastTryingAt:           updated.LastTryingAt,
			Authorized:             updated.Authorized,
		}, nil
	}
	mutex.Lock()
	newRateLimiter, err := rate_limiter_entity.NewRateLimiter(input.IP, input.Token)
	if err != nil {
		log.Println("Erro ao criar novo rate limiter entity:", err)
	}

	err = uc.rateLimiterRepository.InitMonitoring(ctx, newRateLimiter)
	if err != nil {
		log.Println("Erro ao inicializar rate limiter no Redis:", err)
	}
	mutex.Unlock()
	return &InitMonitoringOutputDTO{
		IP:                     newRateLimiter.IP,
		Token:                  newRateLimiter.Token,
		IPLimit:                newRateLimiter.IPLimit,
		TokenLimit:             newRateLimiter.TokenLimit,
		BlockDurationInSeconds: newRateLimiter.BlockDurationInSeconds,
		Reqs:                   newRateLimiter.Reqs,
		InitTryingAt:           newRateLimiter.InitTryingAt,
		LastTryingAt:           newRateLimiter.LastTryingAt,
		Authorized:             newRateLimiter.Authorized,
	}, nil
}

package rate_limiter_usecase

import (
	"context"
	"log"
	"time"

	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_entity"
	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_repository"
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
	InitTryingAt           time.Time `json:"init_trying_at"`
	LastTryingAt           time.Time `json:"last_trying_at"`
	Authorized             bool      `json:"authorized"`
}

func (uc *InitMonitoringUseCase) Execute(ctx context.Context, input InitMonitoringInputDTO) (*InitMonitoringOutputDTO, error) {
	existsLimiter, _ := uc.rateLimiterRepository.FindLimiter(ctx, input.IP, input.Token)
	// if err != nil {
	// 	return nil, err
	// }
	if existsLimiter == nil {
		newRateLimiter, err := rate_limiter_entity.NewRateLimiter(
			input.IP,
			input.Token,
		)
		if err != nil {
			log.Println("error trying create new rate limiter entity")
			return nil, err
		}
		errNew := uc.rateLimiterRepository.InitMonitoring(ctx, newRateLimiter)
		if errNew != nil {
			log.Println("error trying create new rate limiter on repository")
			return nil, errNew
		}

		return &InitMonitoringOutputDTO{
			IP:                     newRateLimiter.IP,
			Token:                  newRateLimiter.Token,
			IPLimit:                newRateLimiter.IPLimit,
			TokenLimit:             newRateLimiter.TokenLimit,
			BlockDurationInSeconds: newRateLimiter.BlockDurationInSeconds,
			InitTryingAt:           newRateLimiter.InitTryingAt,
			LastTryingAt:           newRateLimiter.LastTryingAt,
			Authorized:             newRateLimiter.Authorized,
		}, nil
	}
	// log.Println("identifica update")
	// updateLimiter, err := rate_limiter_entity.UpdateLimiter(*existsLimiter)
	// if err != nil {
	// 	log.Println("error trying update rate limiter entity")
	// 	return nil, err
	// }
	// updateErr := uc.rateLimiterRepository.UpdateMonitoring(ctx, *updateLimiter)
	// if updateErr != nil {
	// 	log.Println("error trying update rate limiter on repository")
	// 	return nil, err
	// }
	// return &InitMonitoringOutputDTO{
	// 	IP:                     updateLimiter.IP,
	// 	Token:                  updateLimiter.Token,
	// 	IPLimit:                updateLimiter.IPLimit,
	// 	TokenLimit:             updateLimiter.TokenLimit,
	// 	BlockDurationInSeconds: updateLimiter.BlockDurationInSeconds,
	// 	InitTryingAt:           updateLimiter.InitTryingAt,
	// 	LastTryingAt:           updateLimiter.LastTryingAt,
	// 	Authorized:             updateLimiter.Authorized,
	// }, nil
	return nil, nil
}

package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_repository"
	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_usecase"
)

type IpTokenRateLimiterMiddleware struct {
	rateLimiterRepository rate_limiter_repository.RateLimiterRepositoryInterface
}

func NewWebIpTokenMiddleware(
	rateLimiterRepository rate_limiter_repository.RateLimiterRepositoryInterface,
) *IpTokenRateLimiterMiddleware {
	return &IpTokenRateLimiterMiddleware{
		rateLimiterRepository: rateLimiterRepository,
	}
}

func (m *IpTokenRateLimiterMiddleware) HttpRateLimiterMiddleware(next http.Handler) http.Handler {
	// escrever contador
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		token := r.Header.Get("API_KEY")

		dto := rate_limiter_usecase.InitMonitoringInputDTO{
			IP:    ip,
			Token: token,
		}

		initMonitoring := rate_limiter_usecase.NewInitMonitoringUseCase(context.Background(), m.rateLimiterRepository)
		output, err := initMonitoring.Execute(context.Background(), dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

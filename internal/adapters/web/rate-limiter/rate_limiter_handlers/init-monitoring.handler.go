package rate_limiter_handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ratelimiter2/internal/rate-limiter/rate_limiter_usecase"
)

func (h *RateLimiterHandler) InitMonitoring(next http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			token := r.Header.Get("API_KEY")
		
			dto := rate_limiter_usecase.InitMonitoringInputDTO{
				IP:    ip,
				Token: token,
			}
			initMonitoring := rate_limiter_usecase.NewInitMonitoringUseCase(context.Background(), h.RateLimiterRepository)
			output, err := initMonitoring.Execute(context.Background(), dto)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		
			if output.Authorized == false {
				log.Println("middleware should block this request: ", output.Authorized)
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(output); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		})	
}

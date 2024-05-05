package rate_limiter_handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"ratelimiter2/internal/rate-limiter/rate_limiter_usecase"
)

func (h *RateLimiterHandler) FindLimiter(w http.ResponseWriter, r *http.Request) {
	findUseCase := rate_limiter_usecase.NewFindLimiterUseCase(context.Background(), h.RateLimiterRepository)
	dto, err := findUseCase.Execute(context.Background(), r.RemoteAddr, r.Header.Get("API_KEY"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

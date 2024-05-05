package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/rate-limiter/rate_limiter_usecase"
)

func (m *IpTokenRateLimiterMiddleware) FindLimiterByKey(w http.ResponseWriter, r *http.Request) {
	log.Println("bunscando", r.RemoteAddr, r.Header.Get("API_KEY"))
	findUseCase := rate_limiter_usecase.NewFindLimiterUseCase(context.Background(), m.rateLimiterRepository)
	dto, err := findUseCase.Execute(context.Background(), r.RemoteAddr, r.Header.Get("API_KEY"))
	log.Println("respost find", dto)
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

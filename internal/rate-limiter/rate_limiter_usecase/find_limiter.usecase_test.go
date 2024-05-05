package rate_limiter_usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ratelimiter2/internal/rate-limiter/rate_limiter_entity"
	"ratelimiter2/internal/rate-limiter/rate_limiter_usecase"

	"github.com/stretchr/testify/assert"
)

func (m *MockRateLimiterRepository) FindLimiter(ctx context.Context, ip, token string) (*rate_limiter_entity.RateLimiter, error) {
	key := ip + ":" + token
	limiter, ok := m.LimiterMap[key]
	if !ok {
		return nil, errors.New("limiter not found")
	}
	return limiter, nil
}

func TestExecute_FindExistingLimiter(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockRateLimiterRepository{
		LimiterMap: make(map[string]*rate_limiter_entity.RateLimiter),
	}

	ip := "192.168.1.1"
	token := "abc123"
	expectedLimiter := &rate_limiter_entity.RateLimiter{
		IP:                     ip,
		Token:                  token,
		IPLimit:                100,
		TokenLimit:             50,
		BlockDurationInSeconds: 60,
		Reqs:                   10,
		InitTryingAt:           time.Now().Add(-10 * time.Second),
		LastTryingAt:           time.Now(),
		Authorized:             true,
	}
	mockRepo.LimiterMap[ip+":"+token] = expectedLimiter

	useCase := rate_limiter_usecase.NewFindLimiterUseCase(ctx, mockRepo)

	output, err := useCase.Execute(ctx, ip, token)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.Equal(t, expectedLimiter.IP, output.IP)
	assert.Equal(t, expectedLimiter.Token, output.Token)
	assert.Equal(t, expectedLimiter.IPLimit, output.IPLimit)
	assert.Equal(t, expectedLimiter.TokenLimit, output.TokenLimit)
	assert.Equal(t, expectedLimiter.BlockDurationInSeconds, output.BlockDurationInSeconds)
	assert.Equal(t, expectedLimiter.Reqs, output.Reqs)
	assert.Equal(t, expectedLimiter.InitTryingAt, output.InitTryingAt)
	assert.Equal(t, expectedLimiter.LastTryingAt, output.LastTryingAt)
	assert.Equal(t, expectedLimiter.Authorized, output.Authorized)
}

func TestExecute_FindNonExistingLimiter(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockRateLimiterRepository{
		LimiterMap: make(map[string]*rate_limiter_entity.RateLimiter),
	}

	useCase := rate_limiter_usecase.NewFindLimiterUseCase(ctx, mockRepo)

	ip := "192.168.1.1"
	token := "abc123"

	output, err := useCase.Execute(ctx, ip, token)
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "limiter not found")
}

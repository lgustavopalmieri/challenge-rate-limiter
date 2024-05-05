package rate_limiter_usecase_test

import (
	"context"
	"os"
	"testing"

	"ratelimiter2/internal/rate-limiter/rate_limiter_entity"
	"ratelimiter2/internal/rate-limiter/rate_limiter_usecase"

	"github.com/stretchr/testify/assert"
)

type MockRateLimiterRepository struct {
	LimiterMap map[string]*rate_limiter_entity.RateLimiter
}

func (m *MockRateLimiterRepository) InitMonitoring(ctx context.Context, limiter *rate_limiter_entity.RateLimiter) error {
	key := limiter.IP + ":" + limiter.Token
	m.LimiterMap[key] = limiter
	return nil
}

func TestExecute(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockRateLimiterRepository{
		LimiterMap: make(map[string]*rate_limiter_entity.RateLimiter),
	}

	useCase := rate_limiter_usecase.NewInitMonitoringUseCase(ctx, mockRepo)

	input := rate_limiter_usecase.InitMonitoringInputDTO{
		IP:    "192.168.1.1",
		Token: "abc123",
	}

	os.Setenv("IP_LIMIT", "100")
	os.Setenv("TOKEN_LIMIT", "50")
	os.Setenv("EXPIRATION_TIME", "60")

	output, err := useCase.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	expectedOutput := &rate_limiter_usecase.InitMonitoringOutputDTO{
		IP:                     "192.168.1.1",
		Token:                  "abc123",
		IPLimit:                100,
		TokenLimit:             50,
		BlockDurationInSeconds: 60,
		Reqs:                   1,
		Authorized:             true,
	}

	assert.Equal(t, expectedOutput.IP, output.IP)
	assert.Equal(t, expectedOutput.Token, output.Token)
	assert.Equal(t, expectedOutput.IPLimit, output.IPLimit)
	assert.Equal(t, expectedOutput.TokenLimit, output.TokenLimit)
	assert.Equal(t, expectedOutput.BlockDurationInSeconds, output.BlockDurationInSeconds)
	assert.Equal(t, expectedOutput.Reqs, output.Reqs)
	assert.Equal(t, expectedOutput.Authorized, output.Authorized)
}

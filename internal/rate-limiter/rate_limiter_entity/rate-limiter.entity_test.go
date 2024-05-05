package rate_limiter_entity_test

import (
	"os"
	"ratelimiter2/internal/rate-limiter/rate_limiter_entity"
	"testing"
	"time"
)

func TestParseEnvToNumber(t *testing.T) {
	tests := []struct {
		envString string
		expected  int64
	}{
		{"123", 123},
		{"-456", -456},
		{"abc", 0},
	}

	for _, test := range tests {
		result := rate_limiter_entity.ParseEnvToNumber(test.envString)
		if result != test.expected {
			t.Errorf("parseEnvToNumber(%s) returned %d, expected %d", test.envString, result, test.expected)
		}
	}
}

func TestRemoveIpPort(t *testing.T) {
	ipWithPort := "192.168.1.1:8080"
	expected := "192.168.1.1"

	result := rate_limiter_entity.RemoveIpPort(ipWithPort)
	if result != expected {
		t.Errorf("RemoveIpPort(%s) returned %s, expected %s", ipWithPort, result, expected)
	}

	ipWithoutPort := "10.0.0.1"
	expected = "10.0.0.1"

	result = rate_limiter_entity.RemoveIpPort(ipWithoutPort)
	if result != expected {
		t.Errorf("RemoveIpPort(%s) returned %s, expected %s", ipWithoutPort, result, expected)
	}
}

func TestNewRateLimiter(t *testing.T) {
	ip := "192.168.1.1"
	token := "abc123"

	os.Setenv("IP_LIMIT", "100")
	os.Setenv("TOKEN_LIMIT", "50")
	os.Setenv("EXPIRATION_TIME", "60")

	rateLimiter, err := rate_limiter_entity.NewRateLimiter(ip, token)
	if err != nil {
		t.Errorf("NewRateLimiter returned an error: %v", err)
	}

	if rateLimiter.IP != "192.168.1.1" {
		t.Errorf("NewRateLimiter IP expected '192.168.1.1', got '%s'", rateLimiter.IP)
	}
	if rateLimiter.Token != "abc123" {
		t.Errorf("NewRateLimiter Token expected 'abc123', got '%s'", rateLimiter.Token)
	}
	if rateLimiter.IPLimit != 100 {
		t.Errorf("NewRateLimiter IPLimit expected 100, got %d", rateLimiter.IPLimit)
	}
}

func TestUpdateLimiter(t *testing.T) {
	rateLimiter := rate_limiter_entity.RateLimiter{
		IP:    "192.168.1.1",
		Token: "abc123",
		Reqs:  10,
	}

	updatedRateLimiter := rate_limiter_entity.UpdateLimiter(rateLimiter)

	if updatedRateLimiter.Reqs != 11 {
		t.Errorf("UpdateLimiter Reqs expected 11, got %d", updatedRateLimiter.Reqs)
	}
	if !updatedRateLimiter.LastTryingAt.After(rateLimiter.LastTryingAt) {
		t.Error("UpdateLimiter LastTryingAt expected to be updated")
	}
}

func TestValidateAuthorize(t *testing.T) {
	rateLimiter := rate_limiter_entity.RateLimiter{
		Token:                  "",
		IPLimit:                100,
		TokenLimit:             50,
		BlockDurationInSeconds: 60,
		Reqs:                   101,
		InitTryingAt:           time.Now().Add(-61 * time.Second),
	}

	authorized := rate_limiter_entity.ValidateAuthorize(rateLimiter)

	if authorized {
		t.Error("ValidateAuthorize expected false, got true")
	}
}

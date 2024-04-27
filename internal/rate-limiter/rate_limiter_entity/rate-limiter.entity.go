package rate_limiter_entity

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type RateLimiter struct {
	IP                     string
	Token                  *string
	IPLimit                int64
	TokenLimit             int64
	BlockDurationInSeconds int64
	InitTryingAt           time.Time
	LastTryingAt           time.Time
	Authorized             bool
}

func parseEnvToNumber(envString string) int64 {
	envNum, err := strconv.ParseInt(envString, 10, 64)
	if err != nil {
		fmt.Printf("error parsing envString to int64: %v\n", err)
		return 0
	}
	return envNum
}

func NewRateLimiter(ip, token string) (*RateLimiter, error) {
	ipLimit := os.Getenv("IP_LIMIT")
	tokenLimit := os.Getenv("TOKEN_LIMIT")
	blockDurationInSeconds := os.Getenv("TOKEN_LIMIT")

	var tokenPtr *string
	if token != "" {
		tokenPtr = &token
	} else {
		token = ""
	}

	return &RateLimiter{
		IP:                     ip,
		Token:                  tokenPtr,
		IPLimit:                parseEnvToNumber(ipLimit),
		TokenLimit:             parseEnvToNumber(tokenLimit),
		BlockDurationInSeconds: parseEnvToNumber(blockDurationInSeconds),
		InitTryingAt:           time.Now(),
		LastTryingAt:           time.Now(),
		Authorized:             true,
	}, nil
}

func UpdateLimiter(rateLimiter RateLimiter) (*RateLimiter, error) {
	return &RateLimiter{
		IP:                     rateLimiter.IP,
		Token:                  rateLimiter.Token,
		IPLimit:                rateLimiter.IPLimit,
		TokenLimit:             rateLimiter.TokenLimit,
		BlockDurationInSeconds: rateLimiter.BlockDurationInSeconds,
		InitTryingAt:           rateLimiter.InitTryingAt,
		LastTryingAt:           time.Now(),
		Authorized:             validateAuthorize(),
	}, nil
}

func validateAuthorize() bool {
	return false
}

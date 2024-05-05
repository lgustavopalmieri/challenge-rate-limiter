package rate_limiter_entity

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type RateLimiter struct {
	IP                     string
	Token                  string
	IPLimit                int64
	TokenLimit             int64
	BlockDurationInSeconds int64
	Reqs                   int64
	InitTryingAt           time.Time
	LastTryingAt           time.Time
	Authorized             bool
}

func ParseEnvToNumber(envString string) int64 {
	envNum, err := strconv.ParseInt(envString, 10, 64)
	if err != nil {
		fmt.Printf("error parsing envString to int64: %v\n", err)
		return 0
	}
	return envNum
}

func RemoveIpPort(receivedIp string) string {
	pos := strings.Index(receivedIp, ":")
	if pos != -1 {
		return receivedIp[:pos]
	}
	return receivedIp
}

func NewRateLimiter(ip, token string) (*RateLimiter, error) {
	ipLimit := os.Getenv("IP_LIMIT")
	tokenLimit := os.Getenv("TOKEN_LIMIT")
	blockDurationInSeconds := os.Getenv("EXPIRATION_TIME")

	var tokenPtr string
	if token != "" {
		tokenPtr = token
	} else {
		token = ""
	}

	return &RateLimiter{
		IP:                     RemoveIpPort(ip),
		Token:                  tokenPtr,
		IPLimit:                ParseEnvToNumber(ipLimit),
		TokenLimit:             ParseEnvToNumber(tokenLimit),
		BlockDurationInSeconds: ParseEnvToNumber(blockDurationInSeconds),
		Reqs:                   1,
		InitTryingAt:           time.Now(),
		LastTryingAt:           time.Now(),
		Authorized:             true,
	}, nil
}

func UpdateLimiter(rateLimiter RateLimiter) *RateLimiter {
	rateLimiter.Reqs++
	updatedRatelimiter := &RateLimiter{
		IP:                     rateLimiter.IP,
		Token:                  rateLimiter.Token,
		IPLimit:                rateLimiter.IPLimit,
		TokenLimit:             rateLimiter.TokenLimit,
		BlockDurationInSeconds: rateLimiter.BlockDurationInSeconds,
		Reqs:                   rateLimiter.Reqs,
		InitTryingAt:           rateLimiter.InitTryingAt,
		LastTryingAt:           time.Now(),
		Authorized:             ValidateAuthorize(rateLimiter),
	}
	return updatedRatelimiter
}

func calculateDiffSeconds(initTryingAt time.Time) int64 {
	calculate := time.Since(initTryingAt)
	difference := calculate.Milliseconds()
	return int64(difference)
}

func ValidateAuthorize(rateLimiter RateLimiter) bool {
	if rateLimiter.Token == "" && rateLimiter.Reqs > rateLimiter.IPLimit {
		return false
	}
	if rateLimiter.Reqs > rateLimiter.TokenLimit {
		return false
	}
	currentRange := calculateDiffSeconds(rateLimiter.InitTryingAt)

	if currentRange < rateLimiter.BlockDurationInSeconds*1000 {
		log.Println("currentRange: ", currentRange, "blockSecs: ", rateLimiter.BlockDurationInSeconds)
		return false
	}

	return true
}

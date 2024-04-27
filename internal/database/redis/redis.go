package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	redisURL := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
		return nil, err
	}
	fmt.Printf("Conexão com o Redis estabelecida. Resposta do PING: %s\n", pong)
	return client, nil
}

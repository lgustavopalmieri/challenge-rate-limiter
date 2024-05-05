package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"ratelimiter2/internal/adapters/web/rate-limiter/rate_limiter_handlers"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Olá! Seu endereço IP é: %s\n", r.RemoteAddr)
}

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}
	redisClient := connectRedis()

	mux := http.NewServeMux()

	handlers := rate_limiter_handlers.NewWebRateLimiterMiddleware(redisClient)

	mux.Handle("/", http.HandlerFunc(helloHandler))

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8686", handlers.InitMonitoring(mux))
}

func connectRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
		return nil
	}
	fmt.Printf("Conexão com o Redis estabelecida. Resposta do PING: %s\n", pong)
	return client
}

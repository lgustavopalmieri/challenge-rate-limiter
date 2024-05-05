package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	middleware "github.com/lgustavopalmieri/challenge-rate-limiter/internal/middleware/http-middleware"
)

// func IPAndTokenMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ip := r.RemoteAddr
// 		token := r.Header.Get("API_KEY")
// 		// if token == "" {
// 		// 	fmt.Printf("Token não fornecido")
// 		// 	http.Error(w, "Token não fornecido", http.StatusUnauthorized)
// 		// 	return
// 		// }

// 		fmt.Printf("Endereço IP do cliente: %s\n", ip)
// 		fmt.Printf("Token da requisição: %s\n", token)

// 		next.ServeHTTP(w, r)
// 	})
// }

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Olá! Seu endereço IP é: %s\n", r.RemoteAddr)
}

func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Adeus! Seu endereço IP é: %s\n", r.RemoteAddr)
}

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	/// fazer di
	redisURL := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
		return
	}
	fmt.Printf("Conexão com o Redis estabelecida. Resposta do PING: %s\n", pong)

	// database.NewRedisClient(context.Background())
	///

	mux := http.NewServeMux()

	routeFind := middleware.NewWebRateLimiterMiddleware(client)

	mux.Handle("/find", http.HandlerFunc(routeFind.FindLimiterByKey))

	mux.Handle("/", http.HandlerFunc(helloHandler))
	mux.Handle("/hello", http.HandlerFunc(helloHandler))
	mux.Handle("/goodbye", http.HandlerFunc(goodbyeHandler))

	fmt.Println("Servidor rodando em http://localhost:8080")

	// globalMiddleware := IPAndTokenMiddleware(mux)
	rateLimiterMiddleware := middleware.NewWebRateLimiterMiddleware(client)

	http.ListenAndServe(":8080", rateLimiterMiddleware.HttpRateLimiterMiddleware(mux))
}

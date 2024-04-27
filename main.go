package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/lgustavopalmieri/challenge-rate-limiter/internal/database/redis"
)

func IPAndTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		fmt.Printf("Endereço IP do cliente: %s\n", ip)

		method := r.Method
		fmt.Printf("Método: %s\n", method)

		token := r.Header.Get("API_KEY")
		if token == "" {
			fmt.Printf("Token não fornecido")
			http.Error(w, "Token não fornecido", http.StatusUnauthorized)
			return
		}

		fmt.Printf("Token da requisição: %s\n", token)

		next.ServeHTTP(w, r)
	})
}

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

	database.NewRedisClient(context.Background())

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(helloHandler))
	mux.Handle("/hello", http.HandlerFunc(helloHandler))
	mux.Handle("/goodbye", http.HandlerFunc(goodbyeHandler))

	fmt.Println("Servidor rodando em http://localhost:8080")
	globalMiddleware := IPAndTokenMiddleware(mux)
	http.ListenAndServe(":8080", globalMiddleware)
}

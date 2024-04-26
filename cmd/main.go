package main

import (
	"fmt"
	"net/http"
)

func IPAndTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		fmt.Printf("Endereço IP do cliente: %s\n", ip)

		route := r.URL.Path
		fmt.Printf("Rota acessada: %s\n", route)

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
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(helloHandler))
	mux.Handle("/hello", http.HandlerFunc(helloHandler))
	mux.Handle("/goodbye", http.HandlerFunc(goodbyeHandler))

	fmt.Println("Servidor rodando em http://localhost:8080")
	globalMiddleware := IPAndTokenMiddleware(mux)
	http.ListenAndServe(":8080", globalMiddleware)
}

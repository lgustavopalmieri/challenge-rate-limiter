## challenge-rate-limiter

Este é o código do desafio de Rate Limiter.

### para testar você pode configurar a variável de ambiente onde:
```
- EXPIRATION_TIME representa os segundos de bloqueio
- IP_LIMIT representa o limite para ips
- TOKEN_LIMIT representa o limite para tokens
```

### para rodar o projeto basta rodar:
```
docker-compose up --build -d
```

### para realizar chamadas:
```
você pode usar uma ferramenta de sua escolha
```

### para substituir o Redis por outro banco de dados:
o rate limiter possui um repositório para firmar um contrato dos métodos necessários
```
type RateLimiterRepositoryInterface interface {
	InitMonitoring(ctx context.Context, rateLimiter *rate_limiter_entity.RateLimiter) error
	FindLimiter(ctx context.Context, ip, token string) (*rate_limiter_entity.RateLimiter, error)
}
```
na pasta internal/database você pode substituir o redis_database por outro banco de dados, implementando os métodos exigidos pelo repository e depois fazer a injeção de dependencia no arquivo: 
```
internal/adapters/web/rate-limiter/rate_limiter_handlers/rate_limiter.di.go
```
e incluir o client de banco de dados escolhido e injetar no handler:
```
func NewWebRateLimiterMiddleware(SEU_CLIENTClient *SEU_CLIENT.Client) *RateLimiterHandler {
	rateLimiterRepository := SEU_CLIENT_database.NewSEU_CLIENTRepositoryDb(context.Background(), SEU_CLIENTClient)
	return NewWebRateLimiterHandlers(rateLimiterRepository)
}
```


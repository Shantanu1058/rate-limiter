package main

import (
	"log"
	"net/http"

	"rate-limiter/limiter"
	"rate-limiter/middleware"
	"rate-limiter/policy"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := policy.LoadPolicyConfig("policy.yaml")
	if err != nil {
		log.Fatal(err)
	}
	resolver := policy.NewExactMatchResolver(cfg)

	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	rLimiter := limiter.NewRedisLimiter(client)

	rlMiddleware := &middleware.RateLimiterMiddleware{
		Resolver: resolver,
		Limiter:  rLimiter,
		IdentityHeader: "X-Api-Key", // optional
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Login success"))
	})

	handler := rlMiddleware.Handler(mux)

	log.Println("Starting server :8080")
	http.ListenAndServe(":8080", handler)
}

package middleware

import (
	"fmt"
	"net/http"

	"rate-limiter/limiter"
	"rate-limiter/policy"
)

type RateLimiterMiddleware struct {
	Resolver       *policy.ExactMatchResolver
	Limiter        *limiter.RedisLimiter
	IdentityHeader string // optional, e.g., "X-Api-Key"
}

func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limits, ok := m.Resolver.Resolve(r.Method, r.URL.Path)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		identity := r.RemoteAddr
		if m.IdentityHeader != "" {
			if id := r.Header.Get(m.IdentityHeader); id != "" {
				identity = id
			}
		}

		// --- Global ---
		globalKey := limiter.BuildGlobalKey(r.Method, r.URL.Path)
		globalLimit := limiter.RateLimit{
			Capacity:       limits.Global.Capacity,
			LeakRatePerSec: limits.Global.LeakRatePerSec,
		}

		globalAllowed, globalRemaining, retryAfter, err := m.Limiter.AllowWithHeaders(globalKey, globalLimit)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", globalLimit.Capacity))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", globalRemaining))
		if !globalAllowed {
			w.Header().Set("X-RateLimit-Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
			http.Error(w, "Too Many Requests (global)", http.StatusTooManyRequests)
			return
		}

		// --- Identity ---
		if limits.Identity != nil {
			identityKey := limiter.BuildIdentityKey(identity, r.Method, r.URL.Path)
			identityLimit := limiter.RateLimit{
				Capacity:       limits.Identity.Capacity,
				LeakRatePerSec: limits.Identity.LeakRatePerSec,
			}

			identityAllowed, identityRemaining, retryAfter, err := m.Limiter.AllowWithHeaders(identityKey, identityLimit)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", identityLimit.Capacity))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", identityRemaining))
			if !identityAllowed {
				w.Header().Set("X-RateLimit-Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
				http.Error(w, "Too Many Requests (identity)", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

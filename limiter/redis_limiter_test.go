package limiter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisLimiter_AllowAndBlock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		t.Fatalf("Redis not reachable: %v", err)
	}

	limiter := NewRedisLimiter(client)
	key := "leaky_bucket:endpoint:POST:/login"
	limit := RateLimit{
		Capacity:       2,
		LeakRatePerSec: 1,
	}

	for i := 1; i <= 3; i++ {
		allowed, err := limiter.Allow(key, limit)
		fmt.Printf("Attempt %d: Allowed=%v, Error=%v\n", i, allowed, err)
		if err != nil {
			t.Fatalf("Redis Eval failed: %v", err)
		}
		if i <= 2 && !allowed {
			t.Fatalf("Attempt %d should have passed", i)
		}
		if i == 3 && allowed {
			t.Fatalf("Attempt %d should have been blocked", i)
		}
	}

	time.Sleep(2 * time.Second)

	allowed, err := limiter.Allow(key, limit)
	if err != nil || !allowed {
		t.Fatalf("Request after leak should have passed: Allowed=%v, Err=%v", allowed, err)
	}
}

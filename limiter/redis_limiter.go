package limiter

import (
	"context"
	_ "embed"
	"strconv"
	"time"
	"github.com/redis/go-redis/v9"
)

//go:embed leaky_bucket.lua
var luaScript string

type RedisLimiter struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisLimiter) Allow(key string, limit RateLimit) (bool, error) {
	now := time.Now().UnixMilli()

	res, err := r.client.Eval(
		r.ctx,
		luaScript,
		[]string{key},
		limit.Capacity,
		limit.LeakRatePerSec,
		now,
	).Int()

	if err != nil {
		return false, err
	}

	return res == 1, nil
}

func (r *RedisLimiter) AllowWithHeaders(key string, limit RateLimit) (allowed bool, remaining int, retryAfter time.Duration, err error) {
	now := time.Now().UnixMilli()

	res, err := r.client.Eval(
		r.ctx,
		luaScript,
		[]string{key},
		limit.Capacity,
		limit.LeakRatePerSec,
		now,
	).Int()
	if err != nil {
		return false, 0, 0, err
	}
	allowed = res == 1

	// Get current water
	waterVal, err := r.client.HGet(r.ctx, key, "water").Result()
	if err != nil && err != redis.Nil {
		return false, 0, 0, err
	}

	water := 0.0
	if waterVal != "" {
		water, _ = strconv.ParseFloat(waterVal, 64)
	}

	// Remaining requests
	remaining = int(float64(limit.Capacity) - water)
	if remaining < 0 {
		remaining = 0
	}

	// Retry after in seconds
	retryAfter = time.Duration((water - float64(limit.Capacity)) / limit.LeakRatePerSec * float64(time.Second))
	if retryAfter < 0 {
		retryAfter = 0
	}

	return allowed, remaining, retryAfter, nil
}
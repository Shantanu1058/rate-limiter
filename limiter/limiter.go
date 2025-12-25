package limiter

type RateLimit struct {
	Capacity       int
	LeakRatePerSec float64
}

type Limiter interface {
	Allow(key string, limit RateLimit) (bool, error)
}

# Rate Limiter(Go + Redis)

A simple **rate limiter** in Go using **Redis** and the **Leaky Bucket algorithm**.  
Supports **global** and **per-identity** limits with standard rate limit headers.

---

## Features

- Global + identity (IP/API key) limits  
- Configurable via `policy.yaml`  
- Middleware for `net/http` servers  
- Returns headers:  
  - `X-RateLimit-Limit`  
  - `X-RateLimit-Remaining`  
  - `X-RateLimit-Retry-After`

---

## Installation

1. Clone repository:

```bash
git clone https://github.com/Shantanu1058/rate-limiter
cd rate-limiter
docker run -d --name redis-rate-limit -p 6379:6379 redis:7
go mod tidy
```

2. Check the policy.yaml file. For example I have added login endpoint

```
domain: api
descriptors:
  - match:
      method: POST
      path: /login
    limits:
      global:
        capacity: 5
        leak_rate_per_sec: 0.1
      identity:
        capacity: 2
        leak_rate_per_sec: 0.05

```
3. Start the server
```
go run main.go
```

4. Test the endpoints
```
curl -i -X POST http://localhost:8080/login
curl -i -X POST http://localhost:8080/login -H "X-Api-Key: user123"
```

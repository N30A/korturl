package middleware

// https://en.wikipedia.org/wiki/Token_bucket
// https://medium.com/@0xTanzim/understanding-the-token-bucket-algorithm-for-rate-limiting-fccdf80e27ca

import (
	"net/http"
	"sync"
	"time"
)

var DefaultRateLimitConfig = RateLimitConfig{
	Capacity:   25,
	RefillRate: 2.5,
	TTL:        time.Second * 120,
}

type RateLimitConfig struct {
	Capacity   int           // the maximum number of tokens in the bucket
	RefillRate float64       // tokens per second
	TTL        time.Duration // time to live for a client's token bucket
}

type RateLimiter struct {
	config  RateLimitConfig
	clients map[string]*client
	mu      sync.Mutex
}

func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		config:  config,
		clients: make(map[string]*client),
	}

	go func() {
		ticker := time.NewTicker(rl.config.TTL)
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for clientAddr, client := range rl.clients {
		if time.Since(client.lastRequest) > rl.config.TTL {
			delete(rl.clients, clientAddr)
		}
	}
}

func (rl *RateLimiter) getOrSetClient(clientAddr string) *client {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, ok := rl.clients[clientAddr]
	if !ok {
		c = &client{tokenBucket: newTokenBucket(rl.config.Capacity, rl.config.RefillRate)}
		rl.clients[clientAddr] = c
	}

	c.lastRequest = time.Now()
	return c
}

type client struct {
	lastRequest time.Time
	tokenBucket *tokenBucket
}

type tokenBucket struct {
	capacity   int       // the maximum number of tokens in the bucket
	tokens     float64   // the current number of tokens in the bucket
	refillRate float64   // tokens per second
	lastRefill time.Time // the time of the last refill
	mu         sync.Mutex
}

func newTokenBucket(capacity int, refillRate float64) *tokenBucket {
	return &tokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *tokenBucket) consume() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = min(float64(tb.capacity), tb.tokens+elapsed*tb.refillRate)
	tb.lastRefill = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}

	return false
}

func RateLimit(rateLimiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client := rateLimiter.getOrSetClient(r.RemoteAddr)
			if !client.tokenBucket.consume() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimitConfig holds the configuration for rate limiting
type RateLimitConfig struct {
	RequestsPerMinute int
	Window            time.Duration
}

// clientInfo stores rate limiting data for a specific client
type clientInfo struct {
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

// rateLimiter implements token bucket algorithm for rate limiting
type rateLimiter struct {
	config  RateLimitConfig
	clients map[string]*clientInfo
	mu      sync.RWMutex
}

// newRateLimiter creates a new rate limiter instance
func newRateLimiter(config RateLimitConfig) *rateLimiter {
	rl := &rateLimiter{
		config:  config,
		clients: make(map[string]*clientInfo),
	}

	// Start cleanup goroutine to remove stale clients
	go rl.cleanup()

	return rl
}

// cleanup removes stale client entries to prevent memory leaks
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, client := range rl.clients {
			client.mu.Lock()
			// Remove clients that haven't been accessed in 2x the window time
			if now.Sub(client.lastRefill) > rl.config.Window*2 {
				delete(rl.clients, ip)
			}
			client.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// getOrCreateClient gets existing client info or creates new one
func (rl *rateLimiter) getOrCreateClient(ip string) *clientInfo {
	rl.mu.RLock()
	client, exists := rl.clients[ip]
	rl.mu.RUnlock()

	if exists {
		return client
	}

	// Client doesn't exist, create it
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check in case another goroutine created it
	if client, exists := rl.clients[ip]; exists {
		return client
	}

	client = &clientInfo{
		tokens:     rl.config.RequestsPerMinute,
		lastRefill: time.Now(),
	}
	rl.clients[ip] = client

	return client
}

// allow checks if a request should be allowed and updates the token count
func (rl *rateLimiter) allow(ip string) (allowed bool, remaining int, resetTime time.Time) {
	client := rl.getOrCreateClient(ip)

	client.mu.Lock()
	defer client.mu.Unlock()

	now := time.Now()

	// Check if window has passed - if so, refill tokens
	if now.Sub(client.lastRefill) >= rl.config.Window {
		client.tokens = rl.config.RequestsPerMinute
		client.lastRefill = now
	}

	// Calculate reset time
	resetTime = client.lastRefill.Add(rl.config.Window)

	// Check if we have tokens available
	if client.tokens > 0 {
		client.tokens--
		return true, client.tokens, resetTime
	}

	return false, 0, resetTime
}

// extractIP extracts the client IP from the request
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For header (used by proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if host, _, err := net.SplitHostPort(xff); err == nil {
			return host
		}
		return xff
	}

	// Check X-Real-IP header (alternative proxy header)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}

	return r.RemoteAddr
}

// RateLimitMiddleware creates a middleware that limits requests per IP address
func RateLimitMiddleware(config RateLimitConfig) func(http.Handler) http.Handler {
	limiter := newRateLimiter(config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r)

			allowed, remaining, resetTime := limiter.allow(ip)

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			if !allowed {
				retryAfter := time.Until(resetTime).Seconds()
				w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter)))
				http.Error(w, fmt.Sprintf("Rate limit exceeded. Try again in %d seconds.", int(retryAfter)), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

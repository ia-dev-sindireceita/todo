package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestRateLimitMiddleware tests the general rate limiting middleware
func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		requestsPerMinute int
		requests          int
		interval          time.Duration
		expectedAllowed   int
		expectedBlocked   int
	}{
		{
			name:              "all requests within limit",
			requestsPerMinute: 100,
			requests:          50,
			interval:          0,
			expectedAllowed:   50,
			expectedBlocked:   0,
		},
		{
			name:              "requests exceed limit",
			requestsPerMinute: 10,
			requests:          15,
			interval:          0,
			expectedAllowed:   10,
			expectedBlocked:   5,
		},
		{
			name:              "single request",
			requestsPerMinute: 100,
			requests:          1,
			interval:          0,
			expectedAllowed:   1,
			expectedBlocked:   0,
		},
		{
			name:              "exact limit reached",
			requestsPerMinute: 5,
			requests:          5,
			interval:          0,
			expectedAllowed:   5,
			expectedBlocked:   0,
		},
		{
			name:              "one over limit",
			requestsPerMinute: 5,
			requests:          6,
			interval:          0,
			expectedAllowed:   5,
			expectedBlocked:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple handler that always returns 200
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create rate limiter with custom config
			config := RateLimitConfig{
				RequestsPerMinute: tt.requestsPerMinute,
				Window:            time.Minute,
			}
			middleware := RateLimitMiddleware(config)
			wrappedHandler := middleware(handler)

			allowed := 0
			blocked := 0

			// Make requests from the same IP
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "192.168.1.1:12345"
				w := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(w, req)

				if w.Code == http.StatusOK {
					allowed++
				} else if w.Code == http.StatusTooManyRequests {
					blocked++
				}

				if tt.interval > 0 {
					time.Sleep(tt.interval)
				}
			}

			if allowed != tt.expectedAllowed {
				t.Errorf("expected %d allowed requests, got %d", tt.expectedAllowed, allowed)
			}

			if blocked != tt.expectedBlocked {
				t.Errorf("expected %d blocked requests, got %d", tt.expectedBlocked, blocked)
			}
		})
	}
}

// TestRateLimitMiddleware_DifferentIPs tests that different IPs have independent limits
func TestRateLimitMiddleware_DifferentIPs(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	config := RateLimitConfig{
		RequestsPerMinute: 5,
		Window:            time.Minute,
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// Make 5 requests from IP1 (should all succeed)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d from IP1 failed with code %d", i+1, w.Code)
		}
	}

	// Make 5 requests from IP2 (should all succeed - independent counter)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d from IP2 failed with code %d", i+1, w.Code)
		}
	}

	// 6th request from IP1 should fail
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 6th request from IP1 to be blocked, got code %d", w.Code)
	}

	// 6th request from IP2 should also fail
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("expected 6th request from IP2 to be blocked, got code %d", w2.Code)
	}
}

// TestRateLimitMiddleware_Headers tests that rate limit headers are set correctly
func TestRateLimitMiddleware_Headers(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	config := RateLimitConfig{
		RequestsPerMinute: 10,
		Window:            time.Minute,
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// First request
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	// Check headers
	if w.Header().Get("X-RateLimit-Limit") != "10" {
		t.Errorf("expected X-RateLimit-Limit to be 10, got %s", w.Header().Get("X-RateLimit-Limit"))
	}

	if w.Header().Get("X-RateLimit-Remaining") != "9" {
		t.Errorf("expected X-RateLimit-Remaining to be 9, got %s", w.Header().Get("X-RateLimit-Remaining"))
	}

	resetHeader := w.Header().Get("X-RateLimit-Reset")
	if resetHeader == "" {
		t.Error("expected X-RateLimit-Reset header to be set")
	}

	// Make 9 more requests to exhaust the limit
	for i := 0; i < 9; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)
	}

	// 11th request should be blocked
	req11 := httptest.NewRequest("GET", "/test", nil)
	req11.RemoteAddr = "192.168.1.1:12345"
	w11 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w11, req11)

	if w11.Code != http.StatusTooManyRequests {
		t.Errorf("expected 11th request to be blocked, got code %d", w11.Code)
	}

	if w11.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Errorf("expected X-RateLimit-Remaining to be 0 on blocked request, got %s", w11.Header().Get("X-RateLimit-Remaining"))
	}
}

// TestRateLimitMiddleware_Reset tests that counter resets after window expires
func TestRateLimitMiddleware_Reset(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Use a short window for testing
	config := RateLimitConfig{
		RequestsPerMinute: 3,
		Window:            200 * time.Millisecond,
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// Make 3 requests (should all succeed)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d failed with code %d", i+1, w.Code)
		}
	}

	// 4th request should fail
	req4 := httptest.NewRequest("GET", "/test", nil)
	req4.RemoteAddr = "192.168.1.1:12345"
	w4 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w4, req4)

	if w4.Code != http.StatusTooManyRequests {
		t.Errorf("4th request should be blocked, got code %d", w4.Code)
	}

	// Wait for window to reset
	time.Sleep(250 * time.Millisecond)

	// Request after reset should succeed
	req5 := httptest.NewRequest("GET", "/test", nil)
	req5.RemoteAddr = "192.168.1.1:12345"
	w5 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w5, req5)

	if w5.Code != http.StatusOK {
		t.Errorf("request after reset failed with code %d", w5.Code)
	}
}

// TestAuthRateLimitMiddleware tests the stricter rate limit for auth endpoints
func TestAuthRateLimitMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	config := RateLimitConfig{
		RequestsPerMinute: 5,
		Window:            time.Minute,
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// Make 5 requests (should all succeed)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("POST", "/auth/login", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d failed with code %d", i+1, w.Code)
		}
	}

	// 6th request should fail
	req6 := httptest.NewRequest("POST", "/auth/login", nil)
	req6.RemoteAddr = "192.168.1.1:12345"
	w6 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w6, req6)

	if w6.Code != http.StatusTooManyRequests {
		t.Errorf("6th auth request should be blocked, got code %d", w6.Code)
	}
}

// TestRateLimitMiddleware_ErrorMessage tests that blocked requests return appropriate message
func TestRateLimitMiddleware_ErrorMessage(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	config := RateLimitConfig{
		RequestsPerMinute: 1,
		Window:            time.Minute,
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// First request succeeds
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req1)

	// Second request should be blocked
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", w2.Code)
	}

	body := w2.Body.String()
	if body == "" {
		t.Error("expected error message in response body")
	}
}

// TestRateLimitMiddleware_IPSpoofingProtection tests protection against IP spoofing
func TestRateLimitMiddleware_IPSpoofingProtection(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Config without trusted proxies - should ignore X-Forwarded-For headers
	config := RateLimitConfig{
		RequestsPerMinute: 3,
		Window:            time.Minute,
		TrustedProxies:    []string{}, // No trusted proxies
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// Make 3 requests with different X-Forwarded-For values
	// All should count as the same IP (192.168.1.100)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.100:12345" // Same real IP
		req.Header.Set("X-Forwarded-For", "10.0.0."+string(rune('1'+i))) // Different forged IPs
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d should succeed (same real IP)", i+1)
		}
	}

	// 4th request should fail because we've exhausted the limit for the real IP
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	req.Header.Set("X-Forwarded-For", "10.0.0.99") // Different forged IP
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Error("spoofed IP should not bypass rate limiting")
	}
}

// TestRateLimitMiddleware_TrustedProxy tests that trusted proxies can set X-Forwarded-For
func TestRateLimitMiddleware_TrustedProxy(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Config with trusted proxy
	config := RateLimitConfig{
		RequestsPerMinute: 3,
		Window:            time.Minute,
		TrustedProxies:    []string{"127.0.0.1"}, // Trust localhost
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// Make 3 requests from trusted proxy with different client IPs
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345" // Trusted proxy
		req.Header.Set("X-Forwarded-For", "10.0.0."+string(rune('1'+i))) // Different client IPs
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d from trusted proxy should succeed", i+1)
		}
	}

	// 4th request from trusted proxy with a new client IP should succeed
	// because it's a different client
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "10.0.0.99")
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("new client IP from trusted proxy should succeed")
	}

	// But 4 requests from the same client IP should fail
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		req.Header.Set("X-Forwarded-For", "10.0.0.1") // Same client IP
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)
	}

	// 4th request from same client should fail
	req4 := httptest.NewRequest("GET", "/test", nil)
	req4.RemoteAddr = "127.0.0.1:12345"
	req4.Header.Set("X-Forwarded-For", "10.0.0.1")
	w4 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w4, req4)

	if w4.Code != http.StatusTooManyRequests {
		t.Error("4th request from same client IP should be blocked")
	}
}

// TestRateLimitMiddleware_UntrustedProxyIgnored tests that untrusted proxies are ignored
func TestRateLimitMiddleware_UntrustedProxyIgnored(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	config := RateLimitConfig{
		RequestsPerMinute: 3,
		Window:            time.Minute,
		TrustedProxies:    []string{"127.0.0.1"}, // Only trust localhost
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// Make 3 requests from untrusted proxy with different X-Forwarded-For values
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.100:12345" // Untrusted proxy
		req.Header.Set("X-Forwarded-For", "10.0.0."+string(rune('1'+i))) // Different forged IPs
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d should succeed", i+1)
		}
	}

	// 4th request should fail (X-Forwarded-For ignored, same real IP)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	req.Header.Set("X-Forwarded-For", "10.0.0.99")
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Error("untrusted proxy should not be able to set client IP")
	}
}

// TestRateLimitMiddleware_XForwardedForMultipleIPs tests parsing of multiple IPs in X-Forwarded-For
func TestRateLimitMiddleware_XForwardedForMultipleIPs(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	config := RateLimitConfig{
		RequestsPerMinute: 3,
		Window:            time.Minute,
		TrustedProxies:    []string{"127.0.0.1"},
	}
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)

	// Request with multiple IPs in X-Forwarded-For (client, proxy1, proxy2)
	// Should use the first IP (10.0.0.1)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		req.Header.Set("X-Forwarded-For", "10.0.0.1, 172.16.0.1, 192.168.1.1")
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d should succeed", i+1)
		}
	}

	// 4th request should fail (same client IP)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "10.0.0.1, 172.16.0.1, 192.168.1.1")
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Error("4th request from same client should be blocked")
	}

	// Request with different first IP should succeed
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	req2.Header.Set("X-Forwarded-For", "10.0.0.2, 172.16.0.1, 192.168.1.1")
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Error("request from different client should succeed")
	}
}

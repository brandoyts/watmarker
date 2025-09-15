package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/port"
)

func RateLimit(cache port.RateLimit, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)
			defer cancel()

			key := clientKey(r, window)

			count, err := cache.Increment(ctx, key)
			if err != nil {
				w.Header().Set("X-RateLimit-Error", "cache_unavailable")
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				_ = cache.Expire(ctx, key, window)
			}

			remaining := limit - int(count)
			if remaining < 0 {
				remaining = 0
			}

			reset := time.Now().Add(window).Unix()

			// observability headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset, 10))

			if int(count) > limit {
				w.Header().Set("Retry-After", strconv.Itoa(int(window.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientKey(r *http.Request, window time.Duration) string {
	// check for forwarded headers (common in docker/proxy setups)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ip = strings.Split(ip, ",")[0] // take first (client) IP
	} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		ip = realIP
	} else {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip = host
		} else {
			ip = r.RemoteAddr
		}
	}

	now := time.Now().Unix()
	bucket := now / int64(window.Seconds())

	return fmt.Sprintf("rate:%s:%d", ip, bucket)
}

package middleware

import (
	"net/http"
	"simple-crud-rnd/helpers"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

var buckets = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getBucketLimiter(ip string, limit, durationInSec int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := buckets[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Duration(durationInSec)*time.Second), limit)
		buckets[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(limit, durationInSec int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := getBucketLimiter(ip, limit, durationInSec)

			if limiter.Allow() {
				return next(c)
			}
			return helpers.Response(c, http.StatusTooManyRequests, nil, "Permintaan ditolak, coba lagi nanti")
		}
	}
}

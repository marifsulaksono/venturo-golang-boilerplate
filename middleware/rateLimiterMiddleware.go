package middleware

import (
	"net/http"
	"simple-crud-rnd/helpers"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

var leakyBuckets = make(map[string]*rate.Limiter)
var leakyMu sync.Mutex

func getBucketLimiter(ip string, limit, durationInSec int) *rate.Limiter {
	leakyMu.Lock()
	defer leakyMu.Unlock()

	limiter, exists := leakyBuckets[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Duration(durationInSec)*time.Second), limit)
		leakyBuckets[ip] = limiter
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

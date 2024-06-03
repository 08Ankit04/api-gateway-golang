package rate_limit

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
	limit       int
	window      time.Duration
)

// Initialize sets the Redis client and rate limiting parameters
func Initialize(redisAddr string, requestLimit int, timeWindow int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	limit = requestLimit
	window = time.Duration(timeWindow) * time.Second
}

// Middleware is the rate limiting middleware
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		// Check if the limiter for this IP already exists in Redis
		val, err := redisClient.Get(ctx, ip).Result()
		if err == redis.Nil {
			// New IP, create a new limiter
			if err := redisClient.Set(ctx, ip, "1", window).Err(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else if err == nil {
			count, _ := strconv.Atoi(val)
			if count >= limit {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			if err := redisClient.Incr(ctx, ip).Err(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

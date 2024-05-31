package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"io/ioutil"

	"github.com/api-gateway-golang/internal/auth"
	"github.com/api-gateway-golang/internal/logger"
	"github.com/api-gateway-golang/internal/rate_limit"
	"github.com/api-gateway-golang/internal/router"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
	RateLimiting struct {
		RedisAddr    string `yaml:"redis_addr"`
		RequestLimit int    `yaml:"request_limit"`
		TimeWindow   int    `yaml:"time_window"`
	} `yaml:"rate_limiting"`
	Routes []struct {
		Path        string `yaml:"path"`
		Service     string `yaml:"service"`
		ServicePort int    `yaml:"service_port"`
	} `yaml:"routes"`
}

var (
	config      Config
	redisClient *redis.Client
	ctx         = context.Background()
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RateLimiting.RedisAddr,
	})

	auth.Initialize(config.JWT.Secret)
	rate_limit.Initialize(config.RateLimiting.RedisAddr, config.RateLimiting.RequestLimit, config.RateLimiting.TimeWindow)
}

func main() {
	router := mux.NewRouter()

	for _, route := range config.Routes {
		router.Handle(route.Path, logger.Middleware(auth.Middleware(rate_limit.Middleware(proxy(route.Service, route.ServicePort))))).Methods("GET", "POST", "PUT", "DELETE")
	}

	// Initialize the router with the configured routes
	routes := []router.Route{}
	for _, r := range config.Routes {
		routes = append(routes, router.Route{Path: r.Path, Service: r.Service, ServicePort: r.ServicePort})
	}
	r := router.InitializeRouter(routes)

	port := config.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func authenticate(next http.Handler) http.Handler {
	return auth.Middleware(next)
}

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := rate.NewLimiter(rate.Every(time.Duration(config.RateLimiting.TimeWindow)*time.Second), config.RateLimiting.RequestLimit)

		// Check if the limiter for this IP already exists in Redis
		val, err := redisClient.Get(ctx, ip).Result()
		if err == redis.Nil {
			// New IP, create a new limiter
			if err := redisClient.Set(ctx, ip, 0, time.Duration(config.RateLimiting.TimeWindow)*time.Second).Err(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else if err == nil {
			count := int(val[0])
			if count >= config.RateLimiting.RequestLimit {
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

func proxy(service string, port int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("http://%s:%d%s", service, port, r.RequestURI)
		req, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		req.Header = r.Header
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		if _, err := ioutil.ReadAll(resp.Body); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}

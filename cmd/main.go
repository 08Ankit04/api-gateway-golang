package main

import (
	"log"
	"net/http"
	"os"

	"github.com/api-gateway-golang/internal/auth"
	"github.com/api-gateway-golang/internal/rate_limit"
	internalRouter "github.com/api-gateway-golang/internal/router"

	"github.com/joho/godotenv"

	"gopkg.in/yaml.v2"
)

const (
	filePathConfig    = "config.yaml"
	ServerPortDefault = "8080"

	errEnvFileNotFound   = "Error No .env file found"
	errReadingConfigFile = "Error reading config file: %v"
	errParsingConfigFile = "Error parsing config file: %v"
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
	config Config
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(errEnvFileNotFound)
	}

	data, err := os.ReadFile(filePathConfig)
	if err != nil {
		log.Fatalf(errReadingConfigFile, err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf(errParsingConfigFile, err)
	}

	auth.Initialize(config.JWT.Secret)
	rate_limit.Initialize(config.RateLimiting.RedisAddr, config.RateLimiting.RequestLimit, config.RateLimiting.TimeWindow)
}

func main() {
	// Initialize the router with the configured routes

	routes := []internalRouter.Route{}
	for _, r := range config.Routes {
		routes = append(routes, internalRouter.Route{Path: r.Path, Service: r.Service, ServicePort: r.ServicePort})
	}
	r := internalRouter.InitializeRouter(routes)

	port := config.Server.Port
	if port == "" {
		port = ServerPortDefault
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

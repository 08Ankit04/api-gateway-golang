package model

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

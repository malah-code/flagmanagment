package config

import "os"

type Config struct {
	BackendAddr  string
	EnvToken     string
	GRPCPort     string
	HealthPort   string
	UpstreamTLS  bool
	TLSCertFile  string
	TLSKeyFile   string
	LogFormat    string
}

func Load() *Config {
	return &Config{
		BackendAddr: getEnv("FM_PROXY_BACKEND_ADDR", "localhost:9090"),
		EnvToken:    getEnv("FM_PROXY_ENV_TOKEN", ""),
		GRPCPort:    getEnv("FM_PROXY_GRPC_PORT", "9091"),
		HealthPort:  getEnv("FM_PROXY_HEALTH_PORT", "8081"),
		UpstreamTLS: getEnvBool("FM_PROXY_UPSTREAM_TLS", false),
		TLSCertFile: getEnv("FM_PROXY_TLS_CERT", ""),
		TLSKeyFile:  getEnv("FM_PROXY_TLS_KEY", ""),
		LogFormat:   getEnv("FM_PROXY_LOG_FORMAT", "auto"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "true" || val == "1" {
		return true
	}
	if val == "false" || val == "0" {
		return false
	}
	return fallback
}

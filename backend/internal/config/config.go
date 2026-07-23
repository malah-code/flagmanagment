package config

import "os"

type Config struct {
	BackendPort string
	Env         string // "development" or "production"
	LogFormat   string // "auto", "text", or "json"
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	RedisHost   string
	RedisPort   string
}

func Load() *Config {
	return &Config{
		BackendPort: getEnv("FM_BACKEND_PORT", "8080"),
		Env:         getEnv("FM_ENV", "development"),
		LogFormat:   getEnv("FM_LOG_FORMAT", "auto"),
		DBHost:      getEnv("FM_DB_HOST", "localhost"),
		DBPort:      getEnv("FM_DB_PORT", "5432"),
		DBUser:      getEnv("FM_DB_USER", "flagmgmt"),
		DBPassword:  getEnv("FM_DB_PASSWORD", "flagmgmt_dev"),
		DBName:      getEnv("FM_DB_NAME", "flagmanagment"),
		RedisHost:   getEnv("FM_REDIS_HOST", "localhost"),
		RedisPort:   getEnv("FM_REDIS_PORT", "6379"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

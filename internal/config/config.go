package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	DatabaseURL    string
	OTPTTLSec      int
	OTPMaxAttempts int
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("invalid %s, using default %d", key, def)
		return def
	}
	return i
}

func Load() Config {
	cfg := Config{
		Port:           getenvInt("PORT", 8080),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		OTPTTLSec:      getenvInt("OTP_TTL_SECONDS", 300),
		OTPMaxAttempts: getenvInt("OTP_MAX_ATTEMPTS", 3),
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required, e.g. postgres://user:pass@localhost:5432/otps?sslmode=disable")
	}
	return cfg
}

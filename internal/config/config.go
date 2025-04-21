package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type DBConfig struct {
	URL               string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
	ConnectionTimeout time.Duration
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type AuthConfig struct {
	JWTSecret       string
	JWTTTL          time.Duration
	RefreshTokenTTL time.Duration
	PasswordCost    int
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
	Auth   AuthConfig
	Env    string
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			URL:               getEnv("DB_URL", "postgres://chikalet:root@localhost:5432/pvz"),
			MaxOpenConns:      getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:      getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime:   getEnvAsDuration("DB_CONN_MAX_LIFETIME", time.Hour),
			ConnectionTimeout: getEnvAsDuration("DB_CONN_TIMEOUT", 5*time.Second),
		},
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 5*time.Second),
		},
		Auth: AuthConfig{
			JWTSecret:       getEnv("JWT_SECRET", "default-secret-key"),
			JWTTTL:          getEnvAsDuration("JWT_TTL", 24*time.Hour),
			RefreshTokenTTL: getEnvAsDuration("REFRESH_TOKEN_TTL", 720*time.Hour),
			PasswordCost:    getEnvAsInt("PASSWORD_COST", 12),
		},
		Env: getEnv("APP_ENV", "dev"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (db *DBConfig) ConnectionString() string {

	separator := "?"
	if strings.Contains(db.URL, "?") {
		separator = "&"
	}

	return fmt.Sprintf("%s%spool_max_conns=%d&pool_max_conn_idle_time=%s",
		db.URL,
		separator,
		db.MaxOpenConns,
		db.ConnMaxLifetime.String())
}

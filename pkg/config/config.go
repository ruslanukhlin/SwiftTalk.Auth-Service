package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)
var (
	ErrParseJWTExpiresAt       = errors.New("failed to parse JWT_EXPIRES_AT duration")
	ErrParseJWTRefreshExpiresAt = errors.New("failed to parse JWT_REFRESH_EXPIRES_AT duration") 
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	SecretKey            string
	ExpiresAfter         time.Duration
	RefreshExpiresAfter  time.Duration
}

type Config struct {
	Mode     string
	Port     string
	Postgres *PostgresConfig
	JWT      *JWTConfig
}

func LoadConfigFromEnv() *Config {
	_ = godotenv.Load(".env.local")

	expiresAfter, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_AT"))
	if err != nil {
		log.Fatalf("%s: %v", ErrParseJWTExpiresAt.Error(), err)
	}

	refreshExpiresAfter, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRES_AT"))
	if err != nil {
		log.Fatalf("%s: %v", ErrParseJWTRefreshExpiresAt.Error(), err)
	}

	return &Config{
		Mode:     os.Getenv("MODE"),
		Port:     os.Getenv("PORT"),
		Postgres: &PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
		},
		JWT: &JWTConfig{
			SecretKey:           os.Getenv("JWT_SECRET_KEY"),
			ExpiresAfter:        expiresAfter,
			RefreshExpiresAfter: refreshExpiresAfter,
		},
	}
}

func DNS(c *PostgresConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)
}
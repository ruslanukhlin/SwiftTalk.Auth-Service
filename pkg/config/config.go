package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	ErrParseJWTExpiresAt        = errors.New("failed to parse JWT_EXPIRES_AT duration")
	ErrParseJWTRefreshExpiresAt = errors.New("failed to parse JWT_REFRESH_EXPIRES_AT duration")
)

var (
	once sync.Once
	cfg  *Config
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	SecretKey           string
	ExpiresAfter        time.Duration
	RefreshExpiresAfter time.Duration
	Issuer              string
	Audience            string
}

type Config struct {
	Mode     string
	PortGrpc string
	PortHttp string
	Postgres *PostgresConfig
	JWT      *JWTConfig
}

func LoadConfigFromEnv() *Config {
	once.Do(func() {
		_ = godotenv.Load(".env.local")

		expiresAfter, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_AT"))
		if err != nil {
			log.Fatalf("%s: %v", ErrParseJWTExpiresAt.Error(), err)
		}

		refreshExpiresAfter, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRES_AT"))
		if err != nil {
			log.Fatalf("%s: %v", ErrParseJWTRefreshExpiresAt.Error(), err)
		}

		cfg = &Config{
			Mode:     os.Getenv("MODE"),
			PortGrpc: os.Getenv("PORT_GRPC"),
			PortHttp: os.Getenv("PORT_HTTP"),
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
				Issuer:              os.Getenv("JWT_ISSUER"),
				Audience:            os.Getenv("JWT_AUDIENCE"),
			},
		}
	})

	return cfg
}

func DNS(c *PostgresConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)
}

func ParseKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	var err error
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey

	privateKeyBytes, err := os.ReadFile("config/keys/private.pem")
	if err != nil {
		return nil, nil, err
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := os.ReadFile("config/keys/public.pem")
	if err != nil {
		return nil, nil, err
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

package gorm

import (
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/infrastructure/db/postgres"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
)

func Migrate(config *config.Config) error {
	return DB.AutoMigrate(&postgres.User{})
}
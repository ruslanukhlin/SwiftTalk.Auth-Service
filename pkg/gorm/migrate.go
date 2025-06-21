package gorm

import (
	"github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/infrastructure/db/postgres"
	"github.com/ruslanukhlin/SwiftTalk.Auth-service/pkg/config"
)

func Migrate(config *config.Config) error {
	return DB.AutoMigrate(&postgres.User{})
}

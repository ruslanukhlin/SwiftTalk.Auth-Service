package postgres

import (
	domain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
	"gorm.io/gorm"
)

var _ domain.UserRepository = &PostgresMemoryRepository{}

type PostgresMemoryRepository struct {
	db *gorm.DB
}

func NewPostgresMemoryRepository(db *gorm.DB) *PostgresMemoryRepository {
	return &PostgresMemoryRepository{
		db: db,
	}
}

func (r *PostgresMemoryRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}
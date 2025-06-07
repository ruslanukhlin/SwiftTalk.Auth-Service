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

func (r *PostgresMemoryRepository) GetUserByUUID(uuid string) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresMemoryRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
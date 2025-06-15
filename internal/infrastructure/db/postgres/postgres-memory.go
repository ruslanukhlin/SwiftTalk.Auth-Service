package postgres

import (
	"errors"

	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
	passwordDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user/password"
	"gorm.io/gorm"
)

var _ user.UserRepository = &PostgresMemoryRepository{}

type PostgresMemoryRepository struct {
	db *gorm.DB
}

func NewPostgresMemoryRepository(db *gorm.DB) *PostgresMemoryRepository {
	return &PostgresMemoryRepository{
		db: db,
	}
}

func (r *PostgresMemoryRepository) CreateUser(user *user.User) error {
	userDb := &User{
		UUID:      user.UUID,
		Email:     user.Email.Value,
		Password:  user.Password.Value,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return r.db.Create(&userDb).Error
}

func (r *PostgresMemoryRepository) GetUserByUUID(uuid string) (*user.User, error) {
	var userDb User

	if err := r.db.Where("uuid = ?", uuid).First(&userDb).Error; err != nil {
		return nil, err
	}

	user := getUserFromDb(userDb)

	return user, nil
}

func (r *PostgresMemoryRepository) GetUserByEmail(email string) (*user.User, error) {
	var userDb User

	if err := r.db.Where("email = ?", email).First(&userDb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	user := getUserFromDb(userDb)

	return user, nil
}

func (r *PostgresMemoryRepository) IsEmailExists(email string) (bool, error) {
	var userDb User

	if err := r.db.Where("email = ?", email).First(&userDb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func getUserFromDb(userDb User) *user.User {
	return &user.User{
		UUID: userDb.UUID,
		Email: user.Email{
			Value: userDb.Email,
		},
		Password: passwordDomain.HashPassword{
			Value: userDb.Password,
		},
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
	}
}

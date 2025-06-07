package domain

type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) bool
}

type PasswordRepository interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) bool
}
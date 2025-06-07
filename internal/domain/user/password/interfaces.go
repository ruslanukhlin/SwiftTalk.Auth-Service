package password

type PasswordRepository interface {
	HashPassword(password Password) (string, error)
	ComparePassword(password string, hash string) error
}
package user

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByUUID(uuid string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	IsEmailExists(email string) (bool, error)
}
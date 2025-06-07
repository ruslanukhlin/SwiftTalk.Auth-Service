package domain

type UserService interface {
	Register(user *User) error
	Login(email, password string) (*User, error)
	VerifyToken(token string) (*User, error)
	RefreshToken(refreshToken string) (*User, error)
}

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByUUID(uuid string) (*User, error)
	GetUserByEmail(email string) (*User, error)
}
package domain

type UserService interface {
	Register(user *User) error
}

type UserRepository interface {
	CreateUser(user *User) error
}
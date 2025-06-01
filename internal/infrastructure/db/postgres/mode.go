package postgres

import "time"

type User struct {
	UUID      string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string `gorm:"not null;type:varchar(255)"`
	Password  string `gorm:"not null;type:varchar(255)"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}
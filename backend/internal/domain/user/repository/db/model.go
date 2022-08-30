package db

import (
	"time"

	"github.com/rl404/image-randomizer/internal/domain/user/entity"
	"gorm.io/gorm"
)

// User is model for user table.
type User struct {
	ID           int64
	Username     string `gorm:"index:unique_username,unique"`
	PasswordHash string
	PasswordSalt string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (u *User) toEntity() *entity.User {
	return &entity.User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PasswordSalt: u.PasswordSalt,
	}
}

func (db *DB) fromEntity(u entity.User) User {
	return User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PasswordSalt: u.PasswordSalt,
	}
}

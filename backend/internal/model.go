package internal

import "time"

// User is model for user table.
type User struct {
	ID           int
	Username     string
	Password     string
	Token        string
	TokenExpired time.Time
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// Image is model for image table.
type Image struct {
	ID        int
	UserID    int
	Image     string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

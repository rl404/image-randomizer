package entity

// User is entity for user.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	PasswordSalt string
}

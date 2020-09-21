package internal

import "errors"

var (
	// ErrRequiredDB will throw if database address is empty.
	ErrRequiredDB = errors.New("required database address")
	// ErrInvalidDB will throw if database address format is invalid.
	ErrInvalidDB = errors.New("invalid database address")
	// ErrNotFound will throw if data not found.
	ErrNotFound = errors.New("data not found")
	// ErrRequiredUser will throw if username is empty.
	ErrRequiredUser = errors.New("required username")
	// ErrRequiredPass will throw if password is empty.
	ErrRequiredPass = errors.New("required password")
	// ErrUserExist will throw if username is already in db.
	ErrUserExist = errors.New("username already registered")
	// ErrRequiredKey will throw if master key in env is empty.
	ErrRequiredKey = errors.New("required master key")
	// ErrWrongUserPass will throw if wrong username or password when login.
	ErrWrongUserPass = errors.New("wrong username or password")
	// ErrRequiredToken will throw if token is not provided in header.
	ErrRequiredToken = errors.New("required token")
	// ErrInvalidToken will throw if token is invalid.
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken will throw if token is expired.
	ErrExpiredToken = errors.New("expired token")
)

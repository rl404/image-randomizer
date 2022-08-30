package db

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/image-randomizer/internal/domain/user/entity"
	"github.com/rl404/image-randomizer/internal/errors"
	"gorm.io/gorm"
)

// DB contains functions for user database.
type DB struct {
	db *gorm.DB
}

// New to create new user database.
func New(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// GetByUsername to get user by username.
func (db *DB) GetByUsername(ctx context.Context, username string) (*entity.User, int, error) {
	var u User
	if err := db.db.WithContext(ctx).Where("username = ?", username).Take(&u).Error; err != nil {
		if _errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrNotFoundUser, err)
		}
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return u.toEntity(), http.StatusOK, nil
}

// Create to create new user.
func (db *DB) Create(ctx context.Context, data entity.User) (*entity.User, int, error) {
	u := db.fromEntity(data)
	if err := db.db.WithContext(ctx).Create(&u).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return u.toEntity(), http.StatusCreated, nil
}

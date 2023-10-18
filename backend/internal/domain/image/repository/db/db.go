package db

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/image-randomizer/internal/domain/image/entity"
	"github.com/rl404/image-randomizer/internal/errors"
	"gorm.io/gorm"
)

// DB contains functions for image database.
type DB struct {
	db *gorm.DB
}

// New to create new image database.
func New(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// Get to get images.
func (db *DB) Get(ctx context.Context, userID int64) ([]*entity.Image, int, error) {
	var images []Image
	if err := db.db.WithContext(ctx).Where("user_id = ?", userID).Find(&images).Error; err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return db.toEntities(images), http.StatusOK, nil
}

// Create to create new image.
func (db *DB) Create(ctx context.Context, data entity.Image) (*entity.Image, int, error) {
	i := db.fromEntity(data)
	if err := db.db.WithContext(ctx).Create(&i).Error; err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return i.toEntity(), http.StatusCreated, nil
}

// Update to update image.
func (db *DB) Update(ctx context.Context, data entity.Image) (int, error) {
	query := db.db.WithContext(ctx).
		Model(&Image{}).
		Where("id = ? and user_id = ?", data.ID, data.UserID).
		Update("image", data.Image)

	if err := query.Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	if query.RowsAffected == 0 {
		return http.StatusNotFound, stack.Wrap(ctx, errors.ErrNotFoundImage)
	}

	return http.StatusOK, nil
}

// Delete to delete image.
func (db *DB) Delete(ctx context.Context, data entity.Image) (int, error) {
	query := db.db.WithContext(ctx).Where("id = ? and user_id = ?", data.ID, data.UserID).Delete(&Image{})

	if err := query.Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	if query.RowsAffected == 0 {
		return http.StatusNotFound, stack.Wrap(ctx, errors.ErrNotFoundImage)
	}

	return http.StatusOK, nil
}

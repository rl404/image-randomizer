package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/image-randomizer/internal/domain/image/entity"
	"github.com/rl404/image-randomizer/internal/domain/image/repository"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/utils"
)

// Cache contains functions for image cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new image cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// Get to get image.
func (c *Cache) Get(ctx context.Context, userID int64) (data []*entity.Image, code int, err error) {
	key := utils.GetKey("images", "user_id", userID)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.Get(ctx, userID)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// Create to create image.
func (c *Cache) Create(ctx context.Context, data entity.Image) (*entity.Image, int, error) {
	key := utils.GetKey("images", "user_id", data.UserID)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return c.repo.Create(ctx, data)
}

// Update to update image.
func (c *Cache) Update(ctx context.Context, data entity.Image) (int, error) {
	key := utils.GetKey("images", "user_id", data.UserID)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return c.repo.Update(ctx, data)
}

// Delete to delete image.
func (c *Cache) Delete(ctx context.Context, data entity.Image) (int, error) {
	key := utils.GetKey("images", "user_id", data.UserID)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return c.repo.Delete(ctx, data)
}

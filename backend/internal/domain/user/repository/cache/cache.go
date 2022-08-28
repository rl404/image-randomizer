package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/image-randomizer/internal/domain/user/entity"
	"github.com/rl404/image-randomizer/internal/domain/user/repository"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/utils"
)

// Cache contains functions for user cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new user cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// GetByUsername to get user by username.
func (c *Cache) GetByUsername(ctx context.Context, username string) (data *entity.User, code int, err error) {
	key := utils.GetKey("user", "username", username)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// Create to create new user.
func (c *Cache) Create(ctx context.Context, data entity.User) (*entity.User, int, error) {
	return c.repo.Create(ctx, data)
}

package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/image-randomizer/internal/domain/image/entity"
	"github.com/rl404/image-randomizer/internal/domain/image/repository"
	"github.com/rl404/image-randomizer/internal/errors"
)

type client struct {
	http *http.Client
	repo repository.Repository
}

// New to create new http client.
func New(repo repository.Repository) *client {
	return &client{
		http: &http.Client{
			Timeout:   30 * time.Second,
			Transport: newrelic.NewRoundTripper(http.DefaultTransport),
		},
		repo: repo,
	}
}

// Get to get image.
func (c *client) Get(ctx context.Context, userID int64) ([]*entity.Image, int, error) {
	return c.repo.Get(ctx, userID)
}

// Create to create image.
func (c *client) Create(ctx context.Context, data entity.Image) (*entity.Image, int, error) {
	return c.repo.Create(ctx, data)
}

// Update to update image.
func (c *client) Update(ctx context.Context, data entity.Image) (int, error) {
	return c.repo.Update(ctx, data)
}

// Delete to delete image.
func (c *client) Delete(ctx context.Context, data entity.Image) (int, error) {
	return c.repo.Delete(ctx, data)
}

// Download to download image.
func (c *client) Download(ctx context.Context, path string) (io.ReadCloser, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode)), errors.ErrInvalidImage)
	}

	return resp.Body, http.StatusOK, nil
}

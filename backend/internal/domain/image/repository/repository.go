package repository

import (
	"context"

	"github.com/rl404/image-randomizer/internal/domain/image/entity"
)

// Repository contains functions for image domain.
type Repository interface {
	Get(ctx context.Context, userID int64) ([]*entity.Image, int, error)
	Create(ctx context.Context, data entity.Image) (*entity.Image, int, error)
	Update(ctx context.Context, data entity.Image) (int, error)
	Delete(ctx context.Context, data entity.Image) (int, error)
}

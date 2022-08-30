package repository

import (
	"context"

	"github.com/rl404/image-randomizer/internal/domain/user/entity"
)

// Repository contains functions for user domain.
type Repository interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, int, error)
	Create(ctx context.Context, data entity.User) (*entity.User, int, error)
}

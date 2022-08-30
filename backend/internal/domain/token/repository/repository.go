package repository

import (
	"context"

	"github.com/rl404/image-randomizer/internal/domain/token/entity"
)

// Repository contains functions for token domain.
type Repository interface {
	CreateAccessToken(ctx context.Context, data entity.CreateAccessTokenRequest) (*entity.Token, int, error)
	CreateRefreshToken(ctx context.Context, data entity.CreateRefreshTokenRequest) (*entity.Token, int, error)
	Get(ctx context.Context, token string) int64
	Delete(ctx context.Context, token string) (int, error)
}

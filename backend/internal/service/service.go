package service

import (
	"context"
	"io"

	imageRepository "github.com/rl404/image-randomizer/internal/domain/image/repository"
	tokenRepository "github.com/rl404/image-randomizer/internal/domain/token/repository"
	userRepository "github.com/rl404/image-randomizer/internal/domain/user/repository"
)

// Service contains functions for service.
type Service interface {
	ValidateToken(ctx context.Context, uuid string, userID int64) (int, error)
	RefreshToken(ctx context.Context, data JWTClaim) (*Token, int, error)

	Register(ctx context.Context, data RegisterRequest) (*Token, int, error)
	Login(ctx context.Context, data LoginRequest) (*Token, int, error)

	GetImages(ctx context.Context, userID int64) ([]Image, int, error)
	CreateImage(ctx context.Context, data CreateImageRequest) (*Image, int, error)
	UpdateImage(ctx context.Context, data UpdateImageRequest) (int, error)
	DeleteImage(ctx context.Context, data DeleteImageRequest) (int, error)

	GetRandomImage(ctx context.Context, username string) (io.ReadCloser, int, error)
}

type service struct {
	user  userRepository.Repository
	image imageRepository.Repository
	token tokenRepository.Repository
}

// Ne to create new service.
func New(
	user userRepository.Repository,
	image imageRepository.Repository,
	token tokenRepository.Repository,
) Service {
	return &service{
		user:  user,
		image: image,
		token: token,
	}
}

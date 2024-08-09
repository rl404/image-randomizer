package service

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	tokenEntity "github.com/rl404/image-randomizer/internal/domain/token/entity"
	"github.com/rl404/image-randomizer/internal/domain/user/entity"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/utils"
)

// RegisterRequest is register request model.
type RegisterRequest struct {
	Username string `json:"username" validate:"required" mod:"trim,lcase"`
	Password string `json:"password" validate:"required" mod:"trim"`
}

// Register to register user.
func (s *service) Register(ctx context.Context, data RegisterRequest) (*Token, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	// Check duplicate username.
	userTmp, code, err := s.user.GetByUsername(ctx, data.Username)
	if err != nil && code != http.StatusNotFound {
		return nil, code, stack.Wrap(ctx, err)
	}
	if userTmp != nil {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, errors.ErrDuplicateUsername)
	}

	// Prepare password salt.
	salt := utils.GenerateUUID()

	// Create user.
	user, code, err := s.user.Create(ctx, entity.User{
		Username:     data.Username,
		PasswordHash: utils.EncodePassword(data.Password, salt),
		PasswordSalt: salt,
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	// Create access token.
	accessToken, code, err := s.token.CreateAccessToken(ctx, tokenEntity.CreateAccessTokenRequest{
		UserID:     user.ID,
		AccessUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	// Create refresh token.
	refreshToken, code, err := s.token.CreateRefreshToken(ctx, tokenEntity.CreateRefreshTokenRequest{
		UserID:      user.ID,
		RefreshUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return &Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}, http.StatusCreated, nil
}

// LoginRequest is login request model.
type LoginRequest struct {
	Username string `json:"username" validate:"required" mod:"trim,lcase"`
	Password string `json:"password" validate:"required" mod:"trim"`
}

// Login to login user.
func (s *service) Login(ctx context.Context, data LoginRequest) (*Token, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	// Get user.
	user, code, err := s.user.GetByUsername(ctx, data.Username)
	if err != nil {
		if code == http.StatusNotFound {
			return nil, http.StatusNotFound, stack.Wrap(ctx, errors.ErrInvalidLogin)
		}
		return nil, code, stack.Wrap(ctx, err)
	}

	if user.PasswordHash != utils.EncodePassword(data.Password, user.PasswordSalt) {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, errors.ErrInvalidLogin)
	}

	// Create access token.
	accessToken, code, err := s.token.CreateAccessToken(ctx, tokenEntity.CreateAccessTokenRequest{
		UserID:     user.ID,
		AccessUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	// Create refresh token.
	refreshToken, code, err := s.token.CreateRefreshToken(ctx, tokenEntity.CreateRefreshTokenRequest{
		UserID:      user.ID,
		RefreshUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return &Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}, http.StatusOK, nil
}

// GetRandomImage to get random image.
func (s *service) GetRandomImage(ctx context.Context, username string) (io.ReadCloser, int, error) {
	user, code, err := s.user.GetByUsername(ctx, username)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	images, code, err := s.image.Get(ctx, user.ID)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	if len(images) == 0 {
		return nil, http.StatusNotFound, stack.Wrap(ctx, errors.ErrNotFoundImage)
	}

	randIndex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(images))

	img, code, err := s.image.Download(ctx, images[randIndex].Image)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return img, http.StatusOK, nil
}

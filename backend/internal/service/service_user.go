package service

import (
	"context"
	"math/rand"
	"net/http"

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
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	// Check duplicate username.
	userTmp, code, err := s.user.GetByUsername(ctx, data.Username)
	if err != nil && code != http.StatusNotFound {
		return nil, code, errors.Wrap(ctx, err)
	}
	if userTmp != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, errors.ErrDuplicateUsername)
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
		return nil, code, errors.Wrap(ctx, err)
	}

	// Create access token.
	accessToken, code, err := s.token.CreateAccessToken(ctx, tokenEntity.CreateAccessTokenRequest{
		UserID:     user.ID,
		AccessUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	// Create refresh token.
	refreshToken, code, err := s.token.CreateRefreshToken(ctx, tokenEntity.CreateRefreshTokenRequest{
		UserID:      user.ID,
		RefreshUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
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
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	// Get user.
	user, code, err := s.user.GetByUsername(ctx, data.Username)
	if err != nil {
		if code == http.StatusNotFound {
			return nil, http.StatusBadRequest, errors.Wrap(ctx, errors.ErrInvalidLogin)
		}
		return nil, code, errors.Wrap(ctx, err)
	}

	if user.PasswordHash != utils.EncodePassword(data.Password, user.PasswordSalt) {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, errors.ErrInvalidLogin)
	}

	// Create access token.
	accessToken, code, err := s.token.CreateAccessToken(ctx, tokenEntity.CreateAccessTokenRequest{
		UserID:     user.ID,
		AccessUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	// Create refresh token.
	refreshToken, code, err := s.token.CreateRefreshToken(ctx, tokenEntity.CreateRefreshTokenRequest{
		UserID:      user.ID,
		RefreshUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	return &Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}, http.StatusOK, nil
}

// GetRandomImage to get random image.
func (s *service) GetRandomImage(ctx context.Context, username string) (string, int, error) {
	user, code, err := s.user.GetByUsername(ctx, username)
	if err != nil {
		return "", code, errors.Wrap(ctx, err)
	}

	images, code, err := s.image.Get(ctx, user.ID)
	if err != nil {
		return "", code, errors.Wrap(ctx, err)
	}

	if len(images) == 0 {
		return "", http.StatusNotFound, errors.Wrap(ctx, errors.ErrNotFoundImage)
	}

	randIndex := rand.Intn(len(images))

	return images[randIndex].Image, http.StatusOK, nil
}

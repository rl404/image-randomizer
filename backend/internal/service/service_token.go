package service

import (
	"context"
	"net/http"

	tokenEntity "github.com/rl404/image-randomizer/internal/domain/token/entity"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/utils"
)

// Token is access and refresh token.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// JWTClaim is jwt claim.
type JWTClaim struct {
	UserID      int64  `json:"user_id"`
	AccessUUID  string `json:"-"`
	RefreshUUID string `json:"-"`
}

// RefreshToken to refresh token.
func (s *service) RefreshToken(ctx context.Context, data JWTClaim) (*Token, int, error) {
	// Create access token.
	accessToken, code, err := s.token.CreateAccessToken(ctx, tokenEntity.CreateAccessTokenRequest{
		UserID:     data.UserID,
		AccessUUID: utils.GenerateUUID(),
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	return &Token{
		AccessToken: accessToken.AccessToken,
	}, http.StatusOK, nil
}

// ValidateToken to validate jwt token.
func (s *service) ValidateToken(ctx context.Context, uuid string, userID int64) (int, error) {
	value := s.token.Get(ctx, uuid)
	if value != userID {
		return http.StatusUnauthorized, errors.Wrap(ctx, errors.ErrInvalidToken)
	}
	return http.StatusOK, nil
}

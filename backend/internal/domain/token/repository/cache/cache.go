package cache

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/image-randomizer/internal/domain/token/entity"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/utils"
)

type client struct {
	cacher         cache.Cacher
	accessSecret   string
	accessExpired  time.Duration
	refreshSecret  string
	refreshExpired time.Duration
}

// New to create new token cache.
func New(cacher cache.Cacher,
	as string, ae time.Duration,
	rs string, re time.Duration,
) *client {
	return &client{
		cacher:         cacher,
		accessSecret:   as,
		accessExpired:  ae,
		refreshSecret:  rs,
		refreshExpired: re,
	}
}

// CreateAccessToken to create new access token.
func (c *client) CreateAccessToken(ctx context.Context, data entity.CreateAccessTokenRequest) (*entity.Token, int, error) {
	keyAccess := utils.GetKey("token", data.AccessUUID)
	if err := c.cacher.Set(ctx, keyAccess, data.UserID, c.accessExpired); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaim := accessToken.Claims.(jwt.MapClaims)
	accessClaim["authorized"] = true
	accessClaim["access_uuid"] = data.AccessUUID
	accessClaim["user_id"] = data.UserID
	accessClaim["iat"] = time.Now().UTC().Unix()
	accessClaim["exp"] = time.Now().UTC().Add(c.accessExpired).Unix()
	accessTokenStr, err := accessToken.SignedString([]byte(c.accessSecret))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err)
	}

	return &entity.Token{
		AccessToken: accessTokenStr,
	}, http.StatusOK, nil
}

// CreateRefreshToken to create new refresh token.
func (c *client) CreateRefreshToken(ctx context.Context, data entity.CreateRefreshTokenRequest) (*entity.Token, int, error) {
	keyRefresh := utils.GetKey("token", data.RefreshUUID)
	if err := c.cacher.Set(ctx, keyRefresh, data.UserID, c.refreshExpired); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaim := refreshToken.Claims.(jwt.MapClaims)
	refreshClaim["authorized"] = true
	refreshClaim["refresh_uuid"] = data.RefreshUUID
	refreshClaim["user_id"] = data.UserID
	refreshClaim["iat"] = time.Now().UTC().Unix()
	refreshClaim["exp"] = time.Now().UTC().Add(c.refreshExpired).Unix()
	refreshTokenStr, err := refreshToken.SignedString([]byte(c.refreshSecret))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err)
	}

	return &entity.Token{
		RefreshToken: refreshTokenStr,
	}, http.StatusOK, nil
}

// Get to get token from cache.
func (c *client) Get(ctx context.Context, token string) (userID int64) {
	c.cacher.Get(ctx, utils.GetKey("token", token), &userID)
	return
}

// Delete to delete token from cache.
func (c *client) Delete(ctx context.Context, token string) (int, error) {
	if err := c.cacher.Delete(ctx, utils.GetKey("token", token)); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}
	return http.StatusOK, nil
}

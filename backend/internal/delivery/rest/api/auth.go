package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/service"
	"github.com/rl404/image-randomizer/internal/utils"
)

type ctxJWTClaim struct{}

type tokenType int8

const (
	tokenAccess tokenType = iota + 1
	tokenRefresh
)

func (api *API) getJWTFromRequest(r *http.Request) string {
	// From query.
	query := r.URL.Query().Get("jwt")
	if query != "" {
		return query
	}

	// From header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	// From cookie.
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (api *API) jwtAuth(next http.HandlerFunc, tokenTypes ...tokenType) http.HandlerFunc {
	tokenType := tokenAccess
	if len(tokenTypes) > 0 {
		tokenType = tokenTypes[0]
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTokenStr := api.getJWTFromRequest(r)
		if jwtTokenStr == "" {
			utils.ResponseWithJSON(w, http.StatusUnauthorized, nil, errors.Wrap(r.Context(), errors.ErrRequiredToken))
			return
		}

		// Parse jwt.
		jwtToken, code, err := api.parseJWT(r.Context(), jwtTokenStr, tokenType)
		if err != nil {
			utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
			return
		}

		var uuid string
		switch tokenType {
		case tokenRefresh:
			uuid = jwtToken.RefreshUUID
		default:
			uuid = jwtToken.AccessUUID
		}

		// Validate token.
		if code, err := api.service.ValidateToken(r.Context(), uuid, jwtToken.UserID); err != nil {
			utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
			return
		}

		ctx := context.WithValue(r.Context(), ctxJWTClaim{}, jwtToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) parseJWT(ctx context.Context, jwtTokenStr string, tokenType tokenType) (*service.JWTClaim, int, error) {
	token, err := jwt.Parse(jwtTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInternalServer
		}
		switch tokenType {
		case tokenAccess:
			return []byte(api.accessSecret), nil
		case tokenRefresh:
			return []byte(api.refreshSecret), nil
		default:
			return nil, nil
		}
	})
	if err != nil {
		return nil, http.StatusUnauthorized, errors.Wrap(ctx, errors.ErrInvalidToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer)
	}

	if !token.Valid {
		return nil, http.StatusUnauthorized, errors.Wrap(ctx, errors.ErrInvalidToken)
	}

	var t service.JWTClaim
	userID, _ := claims["user_id"].(float64)
	t.UserID = int64(userID)
	t.AccessUUID, _ = claims["access_uuid"].(string)
	t.RefreshUUID, _ = claims["refresh_uuid"].(string)

	return &t, http.StatusOK, nil
}

func (api *API) getJWTClaimFromContext(ctx context.Context) (*service.JWTClaim, int, error) {
	claims, ok := ctx.Value(ctxJWTClaim{}).(*service.JWTClaim)
	if !ok {
		return nil, http.StatusInternalServerError, errors.ErrInternalServer
	}
	return claims, http.StatusOK, nil
}

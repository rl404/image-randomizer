package api

import (
	"net/http"

	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/utils"
)

// @summary Check Token
// @tags Token
// @produce json
// @param Authorization header string true "Bearer jwt.access.token"
// @success 200 {object} utils.Response{data=service.JWTClaim}
// @failure 401 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /token/check [get]
func (api *API) handleTokenCheck(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	utils.ResponseWithJSON(w, code, claims, errors.Wrap(r.Context(), err))
}

// @summary Refresh Token
// @tags Token
// @produce json
// @param Authorization header string true "Bearer jwt.refresh.token"
// @success 200 {object} utils.Response{data=service.Token}
// @failure 401 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /token/refresh [post]
func (api *API) handleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
		return
	}

	token, code, err := api.service.RefreshToken(r.Context(), *claims)
	if err == nil {
		token.RefreshToken = api.getJWTFromRequest(r)
	}

	utils.ResponseWithJSON(w, code, token, errors.Wrap(r.Context(), err))
}

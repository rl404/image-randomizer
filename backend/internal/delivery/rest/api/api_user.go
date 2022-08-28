package api

import (
	"encoding/json"
	"net/http"

	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/service"
	"github.com/rl404/image-randomizer/internal/utils"
)

// @summary Register.
// @tags User
// @produce json
// @param request body service.RegisterRequest true "request body"
// @success 201 {object} utils.Response{data=service.Token}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /register [post]
func (api *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	var request service.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidRequestFormat, err))
		return
	}

	token, code, err := api.service.Register(r.Context(), request)
	utils.ResponseWithJSON(w, code, token, errors.Wrap(r.Context(), err))
}

// @summary Login.
// @tags User
// @produce json
// @param request body service.LoginRequest true "request body"
// @success 200 {object} utils.Response{data=service.Token}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /login [post]
func (api *API) handleLogin(w http.ResponseWriter, r *http.Request) {
	var request service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidRequestFormat, err))
		return
	}

	token, code, err := api.service.Login(r.Context(), request)
	utils.ResponseWithJSON(w, code, token, errors.Wrap(r.Context(), err))
}

package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors/stack"
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
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidRequestFormat))
		return
	}

	token, code, err := api.service.Register(r.Context(), request)
	utils.ResponseWithJSON(w, code, token, stack.Wrap(r.Context(), err))
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
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidRequestFormat))
		return
	}

	token, code, err := api.service.Login(r.Context(), request)
	utils.ResponseWithJSON(w, code, token, stack.Wrap(r.Context(), err))
}

// @summary Get random image.
// @tags User
// @produce json,jpeg
// @success 200
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /user/{username}/image.jpg [get]
func (api *API) handleRandomImage(w http.ResponseWriter, r *http.Request) {
	image, code, err := api.service.GetRandomImage(r.Context(), chi.URLParam(r, "username"))
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
		return
	}

	utils.ResponseWithImage(r.Context(), w, image)
}

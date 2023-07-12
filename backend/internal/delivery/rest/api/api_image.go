package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/service"
	"github.com/rl404/image-randomizer/internal/utils"
)

// @summary Get images.
// @tags Image
// @produce json
// @param Authorization header string true "Bearer jwt.access.token"
// @success 200 {object} utils.Response{data=[]service.Image}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /images [get]
func (api *API) handleGetImages(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
		return
	}

	images, code, err := api.service.GetImages(r.Context(), claims.UserID)
	utils.ResponseWithJSON(w, code, images, errors.Wrap(r.Context(), err))
}

// @summary Create image.
// @tags Image
// @produce json
// @param Authorization header string true "Bearer jwt.access.token"
// @param request body service.CreateImageRequest true "request body"
// @success 201 {object} utils.Response{data=service.Image}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /images [post]
func (api *API) handleCreateImage(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
		return
	}

	var request service.CreateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidRequestFormat, err))
		return
	}

	request.UserID = claims.UserID

	image, code, err := api.service.CreateImage(r.Context(), request)
	utils.ResponseWithJSON(w, code, image, errors.Wrap(r.Context(), err))
}

// @summary Update image.
// @tags Image
// @produce json
// @param Authorization header string true "Bearer jwt.access.token"
// @param request body service.UpdateImageRequest true "request body"
// @success 200 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /images/{image_id} [patch]
func (api *API) handleUpdateImage(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
		return
	}

	var request service.UpdateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidRequestFormat, err))
		return
	}

	imageID, err := strconv.ParseInt(chi.URLParam(r, "image_id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidRequestFormat, err))
		return
	}

	request.UserID = claims.UserID
	request.ImageID = imageID

	code, err = api.service.UpdateImage(r.Context(), request)
	utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
}

// @summary Delete image.
// @tags Image
// @produce json
// @param Authorization header string true "Bearer jwt.access.token"
// @success 200 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /images/{image_id} [delete]
func (api *API) handleDeleteImage(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
		return
	}

	imageID, err := strconv.ParseInt(chi.URLParam(r, "image_id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidRequestFormat, err))
		return
	}

	code, err = api.service.DeleteImage(r.Context(), service.DeleteImageRequest{
		UserID:  claims.UserID,
		ImageID: imageID,
	})

	utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
}

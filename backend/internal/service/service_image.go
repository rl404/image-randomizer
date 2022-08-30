package service

import (
	"context"
	"net/http"

	"github.com/rl404/image-randomizer/internal/domain/image/entity"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/utils"
)

// Image is image model.
type Image struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Image  string `json:"image"`
}

// GetImages to get images.
func (s *service) GetImages(ctx context.Context, userID int64) ([]Image, int, error) {
	images, code, err := s.image.Get(ctx, userID)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]Image, len(images))
	for i, img := range images {
		res[i] = Image{
			ID:     img.ID,
			UserID: img.UserID,
			Image:  img.Image,
		}
	}

	return res, http.StatusOK, nil
}

// CreateImageRequest is create image request model.
type CreateImageRequest struct {
	UserID int64  `json:"-" validate:"required" swaggerignore:"true"`
	Image  string `json:"image" validate:"required,url" mod:"trim"`
}

// CreateImage to create new image.
func (s *service) CreateImage(ctx context.Context, data CreateImageRequest) (*Image, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	img, code, err := s.image.Create(ctx, entity.Image{
		UserID: data.UserID,
		Image:  data.Image,
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	return &Image{
		ID:     img.ID,
		UserID: img.UserID,
		Image:  img.Image,
	}, http.StatusCreated, nil
}

// UpdateImageRequest is update image request model.
type UpdateImageRequest struct {
	UserID  int64  `json:"-" validate:"required" swaggerignore:"true"`
	ImageID int64  `json:"-" validate:"required" swaggerignore:"true"`
	Image   string `json:"image" validate:"required,url" mod:"trim"`
}

// UpdateImage to create new image.
func (s *service) UpdateImage(ctx context.Context, data UpdateImageRequest) (int, error) {
	if err := utils.Validate(&data); err != nil {
		return http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	if code, err := s.image.Update(ctx, entity.Image{
		ID:     data.ImageID,
		UserID: data.UserID,
		Image:  data.Image,
	}); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}

// DeleteImageRequest is update image request model.
type DeleteImageRequest struct {
	UserID  int64 `validate:"required"`
	ImageID int64 `validate:"required"`
}

// DeleteImage to create new image.
func (s *service) DeleteImage(ctx context.Context, data DeleteImageRequest) (int, error) {
	if err := utils.Validate(&data); err != nil {
		return http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	if code, err := s.image.Delete(ctx, entity.Image{
		ID:     data.ImageID,
		UserID: data.UserID,
	}); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}

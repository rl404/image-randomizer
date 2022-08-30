package db

import (
	"time"

	"github.com/rl404/image-randomizer/internal/domain/image/entity"
)

// Image is model for image table.
type Image struct {
	ID        int64 `gorm:"index:index_id_user_id"`
	UserID    int64 `gorm:"index:index_id_user_id;index:index_user_id"`
	Image     string
	CreatedAt time.Time
}

func (i *Image) toEntity() *entity.Image {
	return &entity.Image{
		ID:     i.ID,
		UserID: i.UserID,
		Image:  i.Image,
	}
}

func (db *DB) toEntities(data []Image) []*entity.Image {
	imgs := make([]*entity.Image, len(data))
	for i, img := range data {
		imgs[i] = img.toEntity()
	}
	return imgs
}

func (db *DB) fromEntity(i entity.Image) Image {
	return Image{
		ID:     i.ID,
		UserID: i.UserID,
		Image:  i.Image,
	}
}

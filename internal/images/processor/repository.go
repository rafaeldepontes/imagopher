package processor

import (
	"github.com/google/uuid"
	"github.com/rafaeldepontes/imagopher/internal/images/model"
)

type Repository interface {
	FindImageByID(id *int64) (*model.ImageEntity, error)
	FindImageByUUID(uuid uuid.UUID) (*model.ImageEntity, error)
	FindImages(int, *int64) ([]model.ImageEntity, error)
	UploadImage(img *model.ImageEntity) (*int64, error)
}

package processor

import (
	"github.com/google/uuid"
	"github.com/rafaeldepontes/imagopher/internal/images/model"
)

type Repository interface {
	FindImageByID(id *uint64) (*model.ImageEntity, error)
	FindImageByUUID(uuid uuid.UUID) (*model.ImageEntity, error)
	FindImages(int, *uint64) ([]model.ImageEntity, error)
	UploadImage(img *model.ImageEntity) (*uint64, error)
}

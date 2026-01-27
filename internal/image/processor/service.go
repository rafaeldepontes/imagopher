package processor

import (
	"mime/multipart"

	"github.com/rafaeldepontes/imagopher/internal/image/model"
)

type Service interface {
	FindImageByID(id uint64) (*model.ImageEntity, error)
	FindImageByUUID(uuid string) (*model.ImageEntity, error)
	FindImages() ([]model.ImageEntity, error)
	TransformImage(img *model.TransformReq) error
	UploadImage(handler *multipart.FileHeader) (string, error)
}

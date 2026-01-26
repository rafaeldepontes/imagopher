package processor

import "github.com/rafaeldepontes/imagopher/internal/image/model"

type Service interface {
	FindImageByID(id uint64) (*model.ImageEntity, error)
	FindImageByUUID(uuid string) (*model.ImageEntity, error)
	FindImages() ([]model.ImageEntity, error)
	TransformImage(img *model.TransformReq) error
	UploadImage(img *model.ImageEntity) (uint64, error)
}

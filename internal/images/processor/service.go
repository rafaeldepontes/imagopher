package processor

import (
	"mime/multipart"

	cursorModel "github.com/rafaeldepontes/imagopher/internal/cursor/model"
	imgModel "github.com/rafaeldepontes/imagopher/internal/images/model"
)

type Service interface {
	FindImageByID(id *int64) (*imgModel.ImageEntity, error)
	FindImageByUUID(uuid string) (*imgModel.ImageEntity, error)
	FindImages(cursorReq cursorModel.CursorBody) (*imgModel.ImageResp, error)
	TransformImage(img *imgModel.TransformReq) error
	UploadImage(file multipart.File, handler *multipart.FileHeader) (string, error)
}

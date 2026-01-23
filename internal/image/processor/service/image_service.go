package service

import (
	"github.com/rafaeldepontes/imagopher/internal/image/processor"
	"github.com/rafaeldepontes/imagopher/internal/image/processor/repository"
)

type imageService struct {
	repo processor.Repository
}

func NewService() processor.Service {
	return &imageService{
		repo: repository.NewRepository(),
	}
}

// FindImageByID implements [processor.Service].
func (i *imageService) FindImageByID(id uint64) (any, error) {
	panic("unimplemented")
}

// FindImages implements [processor.Service].
func (i *imageService) FindImages() (any, error) {
	panic("unimplemented")
}

// TransformImage implements [processor.Service].
func (i *imageService) TransformImage(img any) error {
	panic("unimplemented")
}

// UploadImage implements [processor.Service].
func (i *imageService) UploadImage(img any) (uint64, error) {
	panic("unimplemented")
}

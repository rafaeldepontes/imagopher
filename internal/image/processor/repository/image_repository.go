package repository

import (
	"database/sql"

	"github.com/rafaeldepontes/imagopher/internal/image/processor"
	"github.com/rafaeldepontes/imagopher/pkg/database/postgres"
)

// TODO: check how to save the images, a BLOB storage would be
// amazing, but I really don't want to use a S3... maybe I
// could use the Potsgres to storage it using the BLOB object
// but I'm not sure how good is...

type imageRepository struct {
	db *sql.DB
}

func NewRepository() processor.Repository {
	return &imageRepository{
		db: postgres.GetDB(),
	}
}

// FindImageByID implements [processor.Repository].
func (i *imageRepository) FindImageByID(id uint64) (any, error) {
	panic("unimplemented")
}

// FindImages implements [processor.Repository].
func (i *imageRepository) FindImages() (any, error) {
	panic("unimplemented")
}

// TransformImage implements [processor.Repository].
func (i *imageRepository) TransformImage(img any) error {
	panic("unimplemented")
}

// UploadImage implements [processor.Repository].
func (i *imageRepository) UploadImage(img any) (uint64, error) {
	panic("unimplemented")
}

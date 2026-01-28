package service

import (
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/imagopher/internal/cache"
	"github.com/rafaeldepontes/imagopher/internal/cache/image"
	"github.com/rafaeldepontes/imagopher/internal/image/model"
	"github.com/rafaeldepontes/imagopher/internal/image/processor"
	"github.com/rafaeldepontes/imagopher/internal/image/processor/repository"
)

const (
	ImgDir         = "./public/uploads/"
	DefaultPermDir = 0755
)

type imageService struct {
	repo  processor.Repository
	cache cache.Cache[string, uint64]
}

func NewService() processor.Service {
	return &imageService{
		repo:  repository.NewRepository(),
		cache: image.NewCache[uint64](),
	}
}

// FindImageByID implements [processor.Service].
func (i *imageService) FindImageByID(id uint64) (*model.ImageEntity, error) {
	log.Printf("[INFO] Searchig for an image by it's ID: %d\n", id)

	img, err := i.repo.FindImageByID(id)
	if err != nil {
		return nil, err
	}

	// TODO: get the image from disk and sent it over.

	return img, nil
}

func (i *imageService) FindImageByUUID(body string) (*model.ImageEntity, error) {
	if data, has := i.cache.Get(body); has {
		log.Println("[INFO] UUID found in cache, using ID instead")
		return i.FindImageByID(data)
	}

	log.Printf("[INFO] Searchig for an image by it's UUID: %s\n", body)

	imgUUID, err := uuid.Parse(body)
	if err != nil {
		log.Println("[ERROR] Could not parse the string into UUID:", err)
		return nil, err
	}

	img, err := i.repo.FindImageByUUID(imgUUID)
	if err != nil {
		return nil, err
	}

	i.cache.Add(body, img.ID, nil)

	// TODO: get the image from disk and sent it over.

	return img, nil
}

// FindImages implements [processor.Service].
func (i *imageService) FindImages() ([]model.ImageEntity, error) {
	imgs, err := i.repo.FindImages()
	if err != nil {
		return nil, err
	}
	return imgs, nil
}

// TransformImage implements [processor.Service].
func (i *imageService) TransformImage(img *model.TransformReq) error {
	// TODO: Finish this method... this is a bunch of image manipulation
	// should be pretty straight foward, just calculate some vectors or
	// I should use a third party library? Need to put a little bit more
	// thought into it...
	panic("unimplemented")
}

// UploadImage has a little (HUGE) problem, it uses the file name as their kinda "path"
// so, any file can override the same photo if they have the same name, it will not be
// stored at the database, but it will override the existing one. Or if it not override
// it will probably open the wrong one...
//
// My go to solution would be using the UUID name to create a new directory and inside
// of this new directory create a new file with the same name. That MAYBE be the solution.
func (i *imageService) UploadImage(handler *multipart.FileHeader) (string, error) {
	UUID := uuid.New()
	name := path.Join(ImgDir, UUID.String(), handler.Filename)

	img := &model.ImageEntity{
		UUID: UUID,
		Path: name,
	}

	if err := os.Mkdir(name, DefaultPermDir); err != nil {
		log.Printf("[ERROR] Could not create the directory %s, because: %s\n", name, err.Error())
		return "", err
	}

	if err := os.WriteFile(name, nil, 1); err != nil {
		log.Printf("[ERROR] Could not write the file: %s - %s\n", name, err.Error())
		return "", err
	}

	_, err := i.repo.UploadImage(img)
	if err != nil {
		return "", err
	}

	// This would be sweet to have, but unfortunately pgx doesn't have
	// support for LastInsertId... So this line is useless for now.
	// i.cache.Add(img.UUID.String(), id, nil)

	return img.UUID.String(), err
}

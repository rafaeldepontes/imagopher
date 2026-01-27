package service

import (
	"log"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/imagopher/internal/cache"
	"github.com/rafaeldepontes/imagopher/internal/cache/image"
	"github.com/rafaeldepontes/imagopher/internal/image/model"
	"github.com/rafaeldepontes/imagopher/internal/image/processor"
	"github.com/rafaeldepontes/imagopher/internal/image/processor/repository"
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
	return img, nil
}

// TODO: Need to check how to use an UUID for real... I'm kinda lost
// here, because the UUID struct and its Value are different...
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

// TODO: check to see if it's better to use a File here and then do all the
// logic...
func (i *imageService) UploadImage(img *model.ImageEntity) (uint64, error) {
	// TODO: Core method to implement... As I need to create the image file
	// from my client to my server, I basically cannot continue without
	// finishing this...
	panic("unimplemented")
}

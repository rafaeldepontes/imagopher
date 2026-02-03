package service

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/rafaeldepontes/imagopher/internal/cache"
	"github.com/rafaeldepontes/imagopher/internal/cache/imagec"
	"github.com/rafaeldepontes/imagopher/internal/images/model"
	"github.com/rafaeldepontes/imagopher/internal/images/processor"
	"github.com/rafaeldepontes/imagopher/internal/images/processor/repository"
)

type imgType string

const (
	// Paths
	ImgDir = "./public/uploads/"

	// Perms
	DefaultPermDir = 0755

	// Types
	JPG  imgType = "jpg"
	APNG imgType = "apng"
	PNG  imgType = "png"
	GIF  imgType = "gif"
	AVIF imgType = "avif"
	SVG  imgType = "svg"
	WebP imgType = "webp"
)

type imageService struct {
	repo  processor.Repository
	cache cache.Cache[string, uint64]
}

func NewService() processor.Service {
	return &imageService{
		repo:  repository.NewRepository(),
		cache: imagec.NewCache[uint64](),
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
func (i *imageService) TransformImage(transform *model.TransformReq) error {
	imgEntity, err := i.FindImageByUUID(transform.UUID.String())
	if err != nil {
		return err
	}

	outputPath := imgEntity.Path
	opts := []imaging.EncodeOption{}

	f, err := os.Open(imgEntity.Path)
	if err != nil {
		log.Println("[ERROR] Could not open the image: ", imgEntity.Path, err)
		return err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return errors.New("unsupported or invalid image")
	}

	if transform.Resize != nil {
		img = imaging.Resize(
			img,
			transform.Resize.Width,
			transform.Resize.Height,
			imaging.Lanczos,
		)
	}

	// TODO: need to create a validation for the crop limites...
	if transform.Crop != nil {
		rect := image.Rect(
			transform.Crop.X,
			transform.Crop.Y,
			transform.Crop.X+transform.Crop.Width,
			transform.Crop.Y+transform.Crop.Height,
		)
		img = imaging.Crop(img, rect)
	}

	if transform.Rotate != nil {
		img = imaging.Rotate(img, float64(*transform.Rotate), color.Transparent)
	}

	if transform.Mirror != nil && *transform.Mirror {
		img = imaging.FlipH(img)
	}

	if transform.Filters != nil {
		if transform.Filters.Grayscale {
			img = imaging.Grayscale(img)
		}
		if transform.Filters.Sepia {
			img = imaging.AdjustSaturation(img, -100)
			img = imaging.AdjustContrast(img, 10)
		}
	}

	if transform.Watermark != nil && *transform.Watermark {
		w := imaging.New(200, 40, color.NRGBA{0, 0, 0, 80})

		img = imaging.Overlay(
			img,
			w,
			image.Pt(
				img.Bounds().Dx()-220,
				img.Bounds().Dy()-60,
			),
			1.0,
		)
	}

	if transform.Compress != nil {
		switch *transform.Compress {
		case "low":
			opts = append(opts, imaging.JPEGQuality(40))
		case "medium":
			opts = append(opts, imaging.JPEGQuality(70))
		case "high":
			opts = append(opts, imaging.JPEGQuality(90))
		}
	}

	if transform.Format != nil {
		switch *transform.Format {
		case "png":
			outputPath = replaceExt(outputPath, ".png")
			opts = append(opts, imaging.PNGCompressionLevel(png.BestCompression))
		case "jpeg":
			outputPath = replaceExt(outputPath, ".jpg")
			opts = append(opts, imaging.JPEGQuality(85))
		case "webp":
			outputPath = replaceExt(outputPath, ".webp")
		}
	}

	return imaging.Save(img, outputPath, opts...)
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
	dir := path.Join(ImgDir, UUID.String())

	parts := strings.Split(handler.Filename, ".")
	if len(parts) < 2 {
		log.Println("[ERROR] File without a proper type.")
		return "", errors.New("File without a type...")
	}

	img := &model.ImageEntity{
		UUID: UUID,

		// I believe this should work believing that the last element in my
		// array should be the file type.
		Type: string(getImageType(parts[len(parts)-1])),
	}

	if err := os.MkdirAll(dir, DefaultPermDir); err != nil {
		log.Printf("[ERROR] Could not create the directory %s, because: %s\n", dir, err.Error())
		return "", err
	}

	dstPath := path.Join(dir, handler.Filename)

	img.Path = dstPath

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("[ERROR] Could not file:", err)
		return "", err
	}
	defer dst.Close()

	src, _ := handler.Open()
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		log.Println("[ERROR] Could not copy the image content:", err)
		return "", err
	}

	_, err = i.repo.UploadImage(img)
	if err != nil {
		return "", err
	}

	// This would be sweet to have, but unfortunately pgx doesn't have
	// support for LastInsertId... So this line is useless for now.
	// i.cache.Add(img.UUID.String(), id, nil)

	return img.UUID.String(), err
}

func getImageType(src string) imgType {
	jpegOption := map[string]struct{}{
		"jpg":   struct{}{},
		"jpeg":  struct{}{},
		"jfif":  struct{}{},
		"pjpeg": struct{}{},
		"pjp":   struct{}{},
	}

	switch src {
	case string(APNG):
		return APNG
	case string(AVIF):
		return AVIF
	case string(GIF):
		return GIF
	case string(PNG):
		return PNG
	case string(SVG):
		return SVG
	case string(WebP):
		return WebP
	}
	if _, has := jpegOption[src]; has {
		return JPG
	}
	return ""
}

func replaceExt(path, ext string) string {
	return strings.TrimSuffix(path, filepath.Ext(path)) + ext
}

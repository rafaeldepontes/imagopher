package service

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
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
	"github.com/rafaeldepontes/imagopher/internal/cursor"
	cursorModel "github.com/rafaeldepontes/imagopher/internal/cursor/model"
	pageSvc "github.com/rafaeldepontes/imagopher/internal/cursor/service"
	imgModel "github.com/rafaeldepontes/imagopher/internal/images/model"
	"github.com/rafaeldepontes/imagopher/internal/images/processor"
	"github.com/rafaeldepontes/imagopher/internal/images/processor/repository"
)

type imgType string

const (
	// Paths
	ImgDir = "./private/uploads/"

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

	// Defaults
	DefaultCursorSize = 10
)

type imageService struct {
	repo       processor.Repository
	pagination cursor.Page
	cache      cache.Cache[string, *uint64]
}

func NewService() processor.Service {
	return &imageService{
		repo:       repository.NewRepository(),
		pagination: pageSvc.NewService(),
		cache:      imagec.NewCache[*uint64](),
	}
}

// FindImageByID implements [processor.Service].
func (i *imageService) FindImageByID(id *uint64) (*imgModel.ImageEntity, error) {
	log.Printf("[INFO] Searchig for an image by it's ID: %d\n", id)

	img, err := i.repo.FindImageByID(id)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (i *imageService) FindImageByUUID(body string) (*imgModel.ImageEntity, error) {
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

	i.cache.Add(body, new(img.ID))

	return img, nil
}

// FindImages implements [processor.Service].
func (i *imageService) FindImages(cursorReq cursorModel.CursorBody) (*imgModel.ImageResp, error) {
	if err := validateCursorBody(cursorReq); err != nil {
		log.Println("[INFO] Invalid cursor body:", err)
		return nil, err
	}

	imgs, err := i.repo.FindImages(cursorReq.Size, cursorReq.NextCursor)
	if err != nil {
		return nil, err
	}

	var imgDTOs []imgModel.ImageDTO
	for i := range imgs {
		imgDTOs = append(imgDTOs, imgModel.ImageDTO{UUID: imgs[i].UUID})
	}

	page, err := i.pagination.Encode(
		cursorReq.Size+DefaultCursorSize,
		new(*cursorReq.NextCursor+DefaultCursorSize),
	)
	if err != nil {
		return nil, err
	}

	resp := &imgModel.ImageResp{
		Data: imgDTOs,
		Page: page,
	}
	return resp, nil
}

// TransformImage implements [processor.Service].
func (i *imageService) TransformImage(transform *imgModel.TransformReq) error {
	imgEntity, err := i.FindImageByUUID(transform.UUID)
	if err != nil {
		return err
	}

	outputPath := imgEntity.Path
	opts := []imaging.EncodeOption{}

	img, err := getImage(outputPath)
	if err != nil {
		return err
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
func (i *imageService) UploadImage(file multipart.File, handler *multipart.FileHeader) (string, error) {
	UUID := uuid.New()
	dir := path.Join(ImgDir, UUID.String())

	parts := strings.Split(handler.Filename, ".")
	if len(parts) < 2 {
		log.Println("[ERROR] File without a proper type.")
		return "", errors.New("File without a type...")
	}

	buffer := make([]byte, 512)

	if _, err := file.Read(buffer); err != nil {
		return "", err
	}

	img := &imgModel.ImageEntity{
		UUID:     UUID,
		MimeType: http.DetectContentType(buffer),

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

	id, err := i.repo.UploadImage(img)
	if err != nil {
		return "", err
	}
	i.cache.Add(img.UUID.String(), id)

	return img.UUID.String(), err
}

func getImageType(src string) imgType {
	jpegOption := map[string]struct{}{
		"jpg":   {},
		"jpeg":  {},
		"jfif":  {},
		"pjpeg": {},
		"pjp":   {},
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

func validateCursorBody(body cursorModel.CursorBody) error {
	if body.Size < 1 {
		return errors.New("Invalid cursor size")
	}

	if body.NextCursor == nil {
		return errors.New("Invalid next cursor")
	}

	return nil
}

func getImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("[ERROR] Could not open the image: ", path, err)
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Println("[ERROR] Image not supported somehow...:", err)
		return nil, errors.New("Unsupported or invalid image")
	}
	return img, nil
}

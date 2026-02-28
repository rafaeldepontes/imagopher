package service

import (
	"errors"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	cursorModel "github.com/rafaeldepontes/imagopher/internal/cursor/model"
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
	DefaultCursorID = 10
)

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

func validateCursorBody(i *imageService, cursor string) (*cursorModel.CursorBody, error) {
	body, err := i.pagination.Decode(cursor)
	if err != nil {
		return nil, err
	}

	if body.Size < 1 {
		return nil, errors.New("Invalid cursor size")
	}

	if body.NextCursor == nil {
		return nil, errors.New("Invalid next cursor")
	}

	return body, nil
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

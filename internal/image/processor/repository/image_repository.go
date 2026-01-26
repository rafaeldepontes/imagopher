package repository

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/imagopher/internal/image/model"
	"github.com/rafaeldepontes/imagopher/internal/image/processor"
	"github.com/rafaeldepontes/imagopher/pkg/database/postgres"
)

// TODO: check how to save the images, a BLOB storage would be
// amazing, but I really don't want to use a S3... maybe I
// could use the Potsgres to storage it using the BLOB object
// but I'm not sure how good is...
//
// For now I've decided to store everything in disk, I don't want
// to use Azure also because of the possibility to be DDoSed!

type imageRepository struct {
	db *sql.DB
}

func NewRepository() processor.Repository {
	return &imageRepository{
		db: postgres.GetDB(),
	}
}

// FindImageByID this is the big deal, use this instead of the FindImageByUUID.
func (i *imageRepository) FindImageByID(id uint64) (*model.ImageEntity, error) {
	img := &model.ImageEntity{}

	query := "SELECT id, path, uuid FROM images img WHERE img.id = ?"

	if err := i.db.QueryRow(query, id).Scan(img.ID, img.Path, img.UUID); err != nil {
		log.Println("[ERROR] Could not complete the scan:", err)
		return nil, err
	}
	return img, nil
}

// FindImageByUUID should only be used in case you don't have the REAL ID in cache already...
// This is pretty slow considering a table scan trying to match a string with a string...
func (i *imageRepository) FindImageByUUID(id uuid.UUID) (*model.ImageEntity, error) {
	img := &model.ImageEntity{}

	query := "SELECT id, path, uuid FROM images img WHERE img.uuid = ?"

	// TODO: Check where the cache layer should be... If here or in the layer above...
	if err := i.db.QueryRow(query, id).Scan(img.ID, img.Path, img.UUID); err != nil {
		log.Println("[ERROR] Could not complete the scan:", err)
		return nil, err
	}
	return img, nil
}

// FindImages returns the list of all images in the database, it should have a
// cursor pagination allowing a BLAZING FAST queries to the Database... (Node.Js moment)
func (i *imageRepository) FindImages() ([]model.ImageEntity, error) {
	var imgs []model.ImageEntity

	query := "SELECT id, path, uuid FROM images;"

	rows, err := i.db.Query(query)
	if err != nil {
		log.Println("[ERROR] Could not built the query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img model.ImageEntity

		if err := rows.Scan(&img.ID, &img.Path, &img.UUID); err != nil {
			log.Println("[ERROR] Could not complete the scan:", err)
			return nil, err
		}
		imgs = append(imgs, img)
	}

	if rows.Err() != nil {
		log.Println("[ERROR] Could not interate over images:", err)
		return nil, rows.Err()
	}

	return imgs, nil
}

func (i *imageRepository) UploadImage(img *model.ImageEntity) (uint64, error) {
	query := "INSERT INTO images (path, uuid) VALUES (?, ?)"

	result, err := i.db.Exec(query, img.Path, img.UUID)
	if err != nil {
		log.Println("[ERROR] Could not insert the new image:", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("[ERROR] Could not extract the last inserted ID:", err)
		return 0, err
	}
	return uint64(id), nil
}

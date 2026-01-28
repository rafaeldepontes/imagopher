package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rafaeldepontes/imagopher/internal/image/model"
	"github.com/rafaeldepontes/imagopher/internal/image/processor"
	"github.com/rafaeldepontes/imagopher/internal/image/processor/service"
)

type imageController struct {
	imgSvc processor.Service
}

func NewController() processor.Controller {
	return &imageController{
		imgSvc: service.NewService(),
	}
}

// FindImageById implements [processor.Controller].
func (ic *imageController) FindImageByID(w http.ResponseWriter, r *http.Request) {
	// Unmarshal the request body...
	param := r.URL.Query().Get("id")
	if param == "" {
		log.Println("[ERROR] Id not present")
		http.Error(w, "Id is missing", http.StatusBadRequest)
		return
	}

	// Sends the JSON to the service
	img, err := ic.imgSvc.FindImageByUUID(param)
	if err != nil {
		log.Println("[ERROR] Could not find the image: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: take care of the data here... Should not sent the whole image
	// just the its bytes and some metadata, also the UUID for the frontend/
	// API...
	if err = json.NewEncoder(w).Encode(img); err != nil {
		log.Println("[ERROR] Could not encode JSON: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	// Returns the image and the correct http status...
	w.Header().Set("Content-Type", "application/json")
}

// FindImages implements [processor.Controller].
func (ic *imageController) FindImages(w http.ResponseWriter, r *http.Request) {
	var imgs []model.ImageEntity
	imgs, err := ic.imgSvc.FindImages()
	if err != nil {
		log.Println("[ERROR] Could not find the images: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sends back everything...
	if err = json.NewEncoder(w).Encode(imgs); err != nil {
		log.Println("[ERROR] Could not parse JSON: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	// Create the beatiful response
	w.Header().Set("Content-Type", "application/json")
}

// TransformImage implements [processor.Controller].
func (ic *imageController) TransformImage(w http.ResponseWriter, r *http.Request) {
	// Unmarshal the body

	// Sends it over and hope for the best...
	if err := ic.imgSvc.TransformImage(nil); err != nil {
		log.Println("[ERROR] Could not transform the imagem: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UploadImage implements [processor.Controller].
func (ic *imageController) UploadImage(w http.ResponseWriter, r *http.Request) {
	// This thing only supports up to 10 Mb...
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("[ERROR] Could not parse the multipart file:", err)
		http.Error(w, "Max file size is 10 Mb", http.StatusBadRequest)
		return
	}

	f, handler, err := r.FormFile("image")
	if err != nil {
		log.Println("[ERROR] Could not get the image file: ", err)
		http.Error(w, "Something really bad happened...", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Sends the file over the service to be saved somewhere...
	id, err := ic.imgSvc.UploadImage(handler)
	if err != nil {
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]string{
		"id": id,
	}); err != nil {
		log.Println("[ERROR] Could not parse JSON: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

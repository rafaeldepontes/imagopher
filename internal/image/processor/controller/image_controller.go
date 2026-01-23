package controller

import (
	"encoding/json"
	"log"
	"net/http"

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

	// Sends the JSON to the service
	img, err := ic.imgSvc.FindImageByID(uint64(int64(1)))
	if err != nil {
		log.Println("[ERROR] Could not find the image: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Encode the photo or what have this is...
	if err = json.NewEncoder(w).Encode(img); err != nil {
		log.Println("[ERROR] Could not encode JSON: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	// Returns the image and the correct http status...
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

// FindImages implements [processor.Controller].
func (ic *imageController) FindImages(w http.ResponseWriter, r *http.Request) {
	// Unmarshal the body

	// Sends it to the service layer...
	// NOT IMPLEMENTED YET:
	// var imgs model.Image[]
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
	w.WriteHeader(http.StatusOK)
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
	// Unmarshal the body into a file

	// Sends the file over the service to be saved somewhere...
	id, err := ic.imgSvc.UploadImage(nil)
	if err != nil {
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]uint64{
		"id": id,
	}); err != nil {
		log.Println("[ERROR] Could not parse JSON: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

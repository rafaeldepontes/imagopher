package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	imgModel "github.com/rafaeldepontes/imagopher/internal/images/model"
	"github.com/rafaeldepontes/imagopher/internal/images/processor"
	"github.com/rafaeldepontes/imagopher/internal/images/processor/service"
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
	param := r.PathValue("id")
	if param == "" {
		log.Println("[ERROR] Id not present")
		http.Error(w, "Id is missing", http.StatusBadRequest)
		return
	}

	img, err := ic.imgSvc.FindImageByUUID(param)
	if err != nil {
		log.Println("[ERROR] Could not find the image: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	f, err := os.Open(img.Path)
	if err != nil {
		log.Println("[ERROR] Could not open the image by its path: ", err)
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", img.MimeType)
	w.Header().Set("Cache-Control", "public, max-age=86400") // 24h cache life.

	http.ServeContent(w, r, "", img.CreatedAt, f)
}

// FindImages implements [processor.Controller].
func (ic *imageController) FindImages(w http.ResponseWriter, r *http.Request) {
	var cursor string = r.URL.Query().Get("c")

	imgs, err := ic.imgSvc.FindImages(cursor)
	if err != nil {
		log.Println("[ERROR] Could not find the images: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(imgs.Data) == 0 {
		log.Println("[INFO] No results")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Create the beatiful response
	w.Header().Set("Content-Type", "application/json")

	// Sends back everything...
	if err = json.NewEncoder(w).Encode(*imgs); err != nil {
		log.Println("[ERROR] Could not parse JSON: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}
}

// TransformImage implements [processor.Controller].
func (ic *imageController) TransformImage(w http.ResponseWriter, r *http.Request) {
	// Unmarshal the body
	var transform imgModel.TransformReq
	if err := json.NewDecoder(r.Body).Decode(&transform); err != nil {
		log.Println("[ERROR] Could not decode the json:", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	// Sends it over and hope for the best...
	if err := ic.imgSvc.TransformImage(&transform); err != nil {
		log.Println("[ERROR] Could not transform the imagem: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}
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
		http.Error(w, "Something really bad", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Sends the file over the service to be saved somewhere...
	id, err := ic.imgSvc.UploadImage(f, handler)
	if err != nil {
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(map[string]string{
		"id": id,
	}); err != nil {
		log.Println("[ERROR] Could not parse JSON: ", err)
		http.Error(w, "Something went really bad", http.StatusInternalServerError)
		return
	}
}

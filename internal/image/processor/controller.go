package processor

import "net/http"

type Controller interface {
	FindImageByID(w http.ResponseWriter, r *http.Request)
	FindImages(w http.ResponseWriter, r *http.Request)
	UploadImage(w http.ResponseWriter, r *http.Request)
	TransformImage(w http.ResponseWriter, r *http.Request)
}

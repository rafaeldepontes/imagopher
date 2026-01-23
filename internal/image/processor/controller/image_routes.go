package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/imagopher/internal/image/processor"
)

func MapImageRoutes(r *chi.Mux, c processor.Controller) {
	// TODO: implement the image processing service... Need to think about the logic
	// maybe storing everything in disk is not a good idea, but I don't want to use
	// AWS S3, so... Need to think more about it.

	// Accepts pagination
	r.Get("/images", c.FindImages)

	r.Post("/images", c.UploadImage)
	r.Post("/images/:id/transform", c.TransformImage)
	r.Get("/images/:id", c.FindImageByID)
}

package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/imagopher/internal/images/processor"
)

// TODO: Implement or use rate limite... I'm think of maybe implement
// my onw, but that seems a little to much.
func MapImageRoutes(r *chi.Mux, c processor.Controller) {

	r.Get("/images", c.FindImages)

	r.Post("/images", c.UploadImage)
	r.Post("/images/{id}/transform", c.TransformImage)
	r.Get("/images/{id}", c.FindImageByID)
}

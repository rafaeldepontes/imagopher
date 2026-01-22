package handler

import "github.com/go-chi/chi/v5"

// I'm thinking of building my on Radix Tree to not use the
// chi package and use mine instead... But I don't want it
// now, but I'm planning...
func Routes(r *chi.Mux) {

	// TODO: implement the image processing service... Need to think about the logic
	// maybe storing everything in disk is not a good idea, but I don't want to use
	// AWS S3, so... Need to think more about it.

	r.Post("/images", nil)

	r.Post("/images/:id/transform", nil)

	r.Get("GET /images/:id", nil)

	// Accepts pagination
	r.Get("GET /images", nil)
}

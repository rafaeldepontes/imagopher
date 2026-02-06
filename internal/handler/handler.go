package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/imagopher/internal/application"
	"github.com/rafaeldepontes/imagopher/internal/images/processor/controller"
)

// I'm thinking of building my on Radix Tree to not use the
// chi package and use mine instead... But I don't want it
// now, but I'm planning...
func Routes(r *chi.Mux, app application.Application) {
	controller.MapImageRoutes(r, app.ImageController)
}

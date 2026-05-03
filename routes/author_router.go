package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func AuthorRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.Author{DB: db}
	r.Get("/authors", controller.GetAllAuthor)
	r.Get("/authors/{id}", controller.GetAuthorByID)
	r.Post("/authors", controller.CreateAuthor)
	r.Put("/authors/{id}", controller.UpdateAuthor)
	r.Delete("/authors/{id}", controller.DeleteAuthor)
}

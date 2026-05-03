package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func BookRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.Book{DB: db}
	r.Get("/books", controller.GetAllBook)
	r.Get("/books/{id}", controller.GetBookByID)
	r.Post("/books", controller.CreateBook)
	r.Put("/books/{id}", controller.UpdateBook)
	r.Delete("/books/{id}", controller.DeleteBook)
}

package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func PublisherRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.Publisher{DB: db}
	r.Get("/publishers", controller.GetAllPublisher)
	r.Get("/publishers/{id}", controller.GetPublisherByID)
	r.Post("/publishers", controller.CreatePublisher)
	r.Put("/publishers/{id}", controller.UpdatePublisher)
	r.Delete("/publishers/{id}", controller.DeletePublisher)
}

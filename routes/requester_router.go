package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RequesterRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.Requester{DB: db}
	r.Get("/requesters", controller.GetAllRequester)
	r.Get("/requesters/{id}", controller.GetRequesterByID)
	r.Post("/requesters", controller.CreateRequester)
	r.Put("/requesters/{id}", controller.UpdateRequester)
	r.Delete("/requesters/{id}", controller.DeleteRequester)
}

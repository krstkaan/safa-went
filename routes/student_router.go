package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func StudentRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.Student{DB: db}
	r.Get("/students", controller.GetAllStudent)
	r.Get("/students/{id}", controller.GetStudentByID)
	r.Post("/students", controller.CreateStudent)
	r.Put("/students/{id}", controller.UpdateStudent)
	r.Delete("/students/{id}", controller.DeleteStudent)
}

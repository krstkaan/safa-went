package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ClassroomRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.Classroom{DB: db}
	r.Get("/classrooms", controller.GetAllClassroom)
	r.Get("/classrooms/{id}", controller.GetClassroomByID)
	r.Post("/classrooms", controller.CreateClassroom)
	r.Put("/classrooms/{id}", controller.UpdateClassroom)
	r.Delete("/classrooms/{id}", controller.DeleteClassroom)
}

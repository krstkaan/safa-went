package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func UserRoutes(r *chi.Mux, db *gorm.DB) {

	controller := &controllers.User{DB: db}

	r.Get("/user", controller.GetAllUser)
	r.Get("/user/{id}", controller.GetUserByID)
	r.Post("/user", controller.CreateUser)
	r.Put("/user/{id}", controller.UpdateUser)
	r.Delete("/user/{id}", controller.DeleteUser)
}
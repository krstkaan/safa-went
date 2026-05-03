package routes

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"safa-went/http/controllers"
)

func AuthRoutes(r *chi.Mux, db *gorm.DB) {
    controller := &controllers.Auth{DB: db}

    r.Post("/auth/register", controller.Register)
    r.Post("/auth/login", controller.Login)
    r.Get("/auth/user", controller.Me)
    r.Post("/auth/logout", controller.Logout)
}
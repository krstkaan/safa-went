package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ApproverRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.Approver{DB: db}
	r.Get("/approvers", controller.GetAllApprover)
	r.Get("/approvers/{id}", controller.GetApproverByID)
	r.Post("/approvers", controller.CreateApprover)
	r.Put("/approvers/{id}", controller.UpdateApprover)
	r.Delete("/approvers/{id}", controller.DeleteApprover)
}

package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func PrintRequestRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.PrintRequest{DB: db}
	r.Get("/print-requests", controller.GetAllPrintRequest)
	r.Get("/print-requests/export/by-requester", controller.ExportByRequester)
	r.Get("/print-requests/export/comparison", controller.ExportComparison)
	r.Get("/print-requests/export/all", controller.ExportAll)
	r.Get("/print-requests/{id}", controller.GetPrintRequestByID)
	r.Post("/print-requests", controller.CreatePrintRequest)
	r.Put("/print-requests/{id}", controller.UpdatePrintRequest)
	r.Delete("/print-requests/{id}", controller.DeletePrintRequest)
}

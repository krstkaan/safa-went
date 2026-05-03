package routes

import (
	"safa-went/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func LoanRoutes(r *chi.Mux, db *gorm.DB) {
	controller := &controllers.Loan{DB: db}
	r.Get("/loans", controller.GetAllLoan)
	r.Get("/loans/{id}", controller.GetLoanByID)
	r.Post("/loans", controller.CreateLoan)
	r.Put("/loans/{id}", controller.UpdateLoan)
	r.Delete("/loans/{id}", controller.DeleteLoan)
	r.Post("/loans/{id}/return", controller.ReturnBook)
	r.Post("/loans/check-availability", controller.CheckAvailability)
}

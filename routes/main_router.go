package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/gorm"
)

func SetupRoutes(r *chi.Mux, db *gorm.DB) {

	// Basic Health Check Endpoint
	r.Get("/ping", HealthCheck(db))

	// Swagger setup
	r.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// User routes
	UserRoutes(r, db)
	// Auth routes
	AuthRoutes(r, db)
	// [*RP*] Please do not delete this comment. It is used to automatically add new route files.
}

// HealthCheck godoc
// @Summary Service health status
// @Description Returns API health and database connectivity status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /ping [get]
func HealthCheck(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			render.JSON(w, r, map[string]interface{}{
				"status":   "degraded",
				"service":  "up",
				"database": "down",
				"error":    "database connection is not initialized",
			})
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			render.JSON(w, r, map[string]interface{}{
				"status":   "degraded",
				"service":  "up",
				"database": "down",
				"error":    err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			render.JSON(w, r, map[string]interface{}{
				"status":   "degraded",
				"service":  "up",
				"database": "down",
				"error":    err.Error(),
			})
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"status":   "ok",
			"service":  "up",
			"database": "up",
		})
	}
}

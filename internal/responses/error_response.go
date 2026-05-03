package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorBody is the standard API error response payload.
type ErrorBody struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// JSONError writes a standard error response.
func JSONError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, ErrorBody{
		Success: false,
		Message: message,
	})
}

// JSONErrorWithDetails writes a standard error response with extra details.
func JSONErrorWithDetails(w http.ResponseWriter, r *http.Request, status int, message string, details interface{}) {
	render.Status(r, status)
	render.JSON(w, r, ErrorBody{
		Success: false,
		Message: message,
		Details: details,
	})
}

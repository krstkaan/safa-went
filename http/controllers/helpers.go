package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ParseID parses the {id} URL parameter into uint.
func ParseID(r *http.Request) (uint, error) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

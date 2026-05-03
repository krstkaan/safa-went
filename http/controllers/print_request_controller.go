package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"gorm.io/gorm"

	"safa-went/database/models"
	"safa-went/http/requests"
	"safa-went/http/resources"
	"safa-went/internal/responses"
)

type PrintRequest struct {
	DB *gorm.DB
}

func (c *PrintRequest) findRaw(id uint, item *models.PrintRequest) error {
	return c.DB.Preload("Requester").Preload("Approver").First(item, id).Error
}

// GetAllPrintRequest godoc
// @Summary List all PrintRequest records
// @Description Get a paginated list of print requests with optional filters
// @Tags PrintRequest
// @Produce json
// @Security BearerAuth
// @Param page               query int    false "Page number (default 1)"
// @Param per_page           query int    false "Items per page (default 15)"
// @Param requester_names    query string false "Comma-separated requester names"
// @Param approver_names     query string false "Comma-separated approver names"
// @Param color_copies_min   query int    false "Minimum color copies"
// @Param color_copies_max   query int    false "Maximum color copies"
// @Param bw_copies_min      query int    false "Minimum BW copies"
// @Param bw_copies_max      query int    false "Maximum BW copies"
// @Param requested_at_from  query string false "Start date (RFC3339)"
// @Param requested_at_to    query string false "End date (RFC3339)"
// @Param sort               query string false "Sort column"
// @Success 200 {object} resources.PrintRequestCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /print-requests [get]
func (c *PrintRequest) GetAllPrintRequest(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)

	requesterNames := splitCSV(r.URL.Query().Get("requester_names"))
	approverNames := splitCSV(r.URL.Query().Get("approver_names"))

	colorCopiesMin, _ := strconv.Atoi(r.URL.Query().Get("color_copies_min"))
	colorCopiesMax, _ := strconv.Atoi(r.URL.Query().Get("color_copies_max"))
	bwCopiesMin, _ := strconv.Atoi(r.URL.Query().Get("bw_copies_min"))
	bwCopiesMax, _ := strconv.Atoi(r.URL.Query().Get("bw_copies_max"))

	var requestedAtFrom, requestedAtTo *time.Time
	if v := r.URL.Query().Get("requested_at_from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			requestedAtFrom = &t
		}
	}
	if v := r.URL.Query().Get("requested_at_to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			requestedAtTo = &t
		}
	}

	sortBy := r.URL.Query().Get("sort")

	q := resources.NewPrintRequestQuery(c.DB)
	collection, err := q.Filter(
		requesterNames, approverNames,
		colorCopiesMin, colorCopiesMax,
		bwCopiesMin, bwCopiesMax,
		requestedAtFrom, requestedAtTo,
		sortBy, page, perPage,
	)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, collection)
}

// GetPrintRequestByID godoc
// @Summary Get a PrintRequest by ID
// @Tags PrintRequest
// @Produce json
// @Security BearerAuth
// @Param id path int true "PrintRequest ID"
// @Success 200 {object} resources.PrintRequestResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /print-requests/{id} [get]
func (c *PrintRequest) GetPrintRequestByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewPrintRequestQuery(c.DB).Find(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, item)
}

// CreatePrintRequest godoc
// @Summary Create a PrintRequest
// @Tags PrintRequest
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.PrintRequestPayload true "PrintRequest payload"
// @Success 201 {object} resources.PrintRequestResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /print-requests [post]
func (c *PrintRequest) CreatePrintRequest(w http.ResponseWriter, r *http.Request) {
	var payload requests.PrintRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	item := models.PrintRequest{
		RequestedAt: payload.RequestedAt,
		ColorCopies: payload.ColorCopies,
		BWCopies:    payload.BWCopies,
		Description: payload.Description,
		RequesterID: payload.RequesterID,
		ApproverID:  payload.ApproverID,
	}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	created, err := resources.NewPrintRequestQuery(c.DB).Find(item.ID)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, created)
}

// UpdatePrintRequest godoc
// @Summary Update a PrintRequest
// @Tags PrintRequest
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "PrintRequest ID"
// @Param payload body requests.PrintRequestUpdatePayload true "PrintRequest update payload"
// @Success 200 {object} resources.PrintRequestResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /print-requests/{id} [put]
func (c *PrintRequest) UpdatePrintRequest(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.PrintRequest
	if err := c.findRaw(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.PrintRequestUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	updates := map[string]interface{}{}
	if payload.RequestedAt != nil {
		updates["requested_at"] = *payload.RequestedAt
	}
	if payload.ColorCopies != nil {
		updates["color_copies"] = *payload.ColorCopies
	}
	if payload.BWCopies != nil {
		updates["bw_copies"] = *payload.BWCopies
	}
	if payload.Description != nil {
		updates["description"] = *payload.Description
	}
	if payload.RequesterID != nil {
		updates["requester_id"] = *payload.RequesterID
	}
	if payload.ApproverID != nil {
		updates["approver_id"] = *payload.ApproverID
	}
	if len(updates) == 0 {
		responses.JSONError(w, r, http.StatusBadRequest, "no fields to update")
		return
	}
	if err := c.DB.Model(&item).Updates(updates).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	updated, err := resources.NewPrintRequestQuery(c.DB).Find(id)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, updated)
}

// DeletePrintRequest godoc
// @Summary Delete a PrintRequest
// @Tags PrintRequest
// @Security BearerAuth
// @Param id path int true "PrintRequest ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /print-requests/{id} [delete]
func (c *PrintRequest) DeletePrintRequest(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.PrintRequest
	if err := c.findRaw(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if err := c.DB.Delete(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// splitCSV splits a comma-separated string into a slice of non-empty strings.
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := make([]string, 0)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := s[start:i]
			if part != "" {
				parts = append(parts, part)
			}
			start = i + 1
		}
	}
	return parts
}

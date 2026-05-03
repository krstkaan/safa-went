package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"gorm.io/gorm"

	"safa-went/database/models"
	"safa-went/http/requests"
	"safa-went/http/resources"
	"safa-went/internal/responses"
)

type Approver struct {
	DB *gorm.DB
}

func (c *Approver) findByID(id uint, item *models.Approver) error {
	return c.DB.First(item, id).Error
}

// GetAllApprover godoc
// @Summary List all Approver records
// @Description Get a paginated list of Approver records
// @Tags Approver
// @Produce json
// @Security BearerAuth
// @Param page     query int    false "Page number (default 1)"
// @Param per_page query int    false "Items per page (default 15)"
// @Param search   query string false "Search by name"
// @Success 200 {object} resources.ApproverCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /approvers [get]
func (c *Approver) GetAllApprover(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	search := r.URL.Query().Get("search")

	q := resources.NewApproverQuery(c.DB)
	if search != "" {
		data, err := q.Search(search, page, perPage)
		if err != nil {
			responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		render.JSON(w, r, data)
		return
	}

	collection, err := q.Paginate(page, perPage)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, collection)
}

// GetApproverByID godoc
// @Summary Get an Approver by ID
// @Tags Approver
// @Produce json
// @Security BearerAuth
// @Param id path int true "Approver ID"
// @Success 200 {object} resources.ApproverResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /approvers/{id} [get]
func (c *Approver) GetApproverByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewApproverQuery(c.DB).Find(id)
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

// CreateApprover godoc
// @Summary Create an Approver
// @Tags Approver
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.ApproverPayload true "Approver payload"
// @Success 201 {object} resources.ApproverResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 422 {object} responses.ErrorBody
// @Router /approvers [post]
func (c *Approver) CreateApprover(w http.ResponseWriter, r *http.Request) {
	var payload requests.ApproverPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	if payload.Name == "" {
		responses.JSONError(w, r, http.StatusUnprocessableEntity, "name is required")
		return
	}
	item := models.Approver{Name: payload.Name}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resources.NewApproverResource(item))
}

// UpdateApprover godoc
// @Summary Update an Approver
// @Tags Approver
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Approver ID"
// @Param payload body requests.ApproverUpdatePayload true "Approver update payload"
// @Success 200 {object} resources.ApproverResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /approvers/{id} [put]
func (c *Approver) UpdateApprover(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Approver
	if err := c.findByID(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.ApproverUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	updates := map[string]interface{}{}
	if payload.Name != nil {
		updates["name"] = *payload.Name
	}
	if len(updates) == 0 {
		responses.JSONError(w, r, http.StatusBadRequest, "no fields to update")
		return
	}
	if err := c.DB.Model(&item).Updates(updates).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if err := c.findByID(id, &item); err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, resources.NewApproverResource(item))
}

// DeleteApprover godoc
// @Summary Delete an Approver
// @Tags Approver
// @Security BearerAuth
// @Param id path int true "Approver ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /approvers/{id} [delete]
func (c *Approver) DeleteApprover(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Approver
	if err := c.findByID(id, &item); err != nil {
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

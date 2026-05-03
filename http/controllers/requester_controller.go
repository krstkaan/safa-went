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

type Requester struct {
	DB *gorm.DB
}

func (c *Requester) findByID(id uint, item *models.Requester) error {
	return c.DB.First(item, id).Error
}

// GetAllRequester godoc
// @Summary List all Requester records
// @Description Get a paginated list of Requester records
// @Tags Requester
// @Produce json
// @Security BearerAuth
// @Param page     query int    false "Page number (default 1)"
// @Param per_page query int    false "Items per page (default 15)"
// @Param search   query string false "Search by name"
// @Success 200 {object} resources.RequesterCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /requesters [get]
func (c *Requester) GetAllRequester(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	search := r.URL.Query().Get("search")

	q := resources.NewRequesterQuery(c.DB)
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

// GetRequesterByID godoc
// @Summary Get a Requester by ID
// @Tags Requester
// @Produce json
// @Security BearerAuth
// @Param id path int true "Requester ID"
// @Success 200 {object} resources.RequesterResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /requesters/{id} [get]
func (c *Requester) GetRequesterByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewRequesterQuery(c.DB).Find(id)
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

// CreateRequester godoc
// @Summary Create a Requester
// @Tags Requester
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.RequesterPayload true "Requester payload"
// @Success 201 {object} resources.RequesterResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 422 {object} responses.ErrorBody
// @Router /requesters [post]
func (c *Requester) CreateRequester(w http.ResponseWriter, r *http.Request) {
	var payload requests.RequesterPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	if payload.Name == "" {
		responses.JSONError(w, r, http.StatusUnprocessableEntity, "name is required")
		return
	}
	item := models.Requester{Name: payload.Name}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resources.NewRequesterResource(item))
}

// UpdateRequester godoc
// @Summary Update a Requester
// @Tags Requester
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Requester ID"
// @Param payload body requests.RequesterUpdatePayload true "Requester update payload"
// @Success 200 {object} resources.RequesterResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /requesters/{id} [put]
func (c *Requester) UpdateRequester(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Requester
	if err := c.findByID(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.RequesterUpdatePayload
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
	render.JSON(w, r, resources.NewRequesterResource(item))
}

// DeleteRequester godoc
// @Summary Delete a Requester
// @Tags Requester
// @Security BearerAuth
// @Param id path int true "Requester ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /requesters/{id} [delete]
func (c *Requester) DeleteRequester(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Requester
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

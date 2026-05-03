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

type Publisher struct {
	DB *gorm.DB
}

func (c *Publisher) findByID(id uint, item *models.Publisher) error {
	return c.DB.First(item, id).Error
}

// GetAllPublisher godoc
// @Summary List all Publisher records
// @Description Get a paginated list of Publisher records
// @Tags Publisher
// @Produce json
// @Security BearerAuth
// @Param page     query int    false "Page number (default 1)"
// @Param per_page query int    false "Items per page (default 15)"
// @Param search   query string false "Search by name"
// @Success 200 {object} resources.PublisherCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /publishers [get]
func (c *Publisher) GetAllPublisher(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	search := r.URL.Query().Get("search")

	q := resources.NewPublisherQuery(c.DB)
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

// GetPublisherByID godoc
// @Summary Get a Publisher by ID
// @Tags Publisher
// @Produce json
// @Security BearerAuth
// @Param id path int true "Publisher ID"
// @Success 200 {object} resources.PublisherResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /publishers/{id} [get]
func (c *Publisher) GetPublisherByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewPublisherQuery(c.DB).Find(id)
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

// CreatePublisher godoc
// @Summary Create a Publisher
// @Tags Publisher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.PublisherPayload true "Publisher payload"
// @Success 201 {object} resources.PublisherResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /publishers [post]
func (c *Publisher) CreatePublisher(w http.ResponseWriter, r *http.Request) {
	var payload requests.PublisherPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	item := models.Publisher{Name: payload.Name}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resources.NewPublisherResource(item))
}

// UpdatePublisher godoc
// @Summary Update a Publisher
// @Tags Publisher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Publisher ID"
// @Param payload body requests.PublisherUpdatePayload true "Publisher update payload"
// @Success 200 {object} resources.PublisherResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /publishers/{id} [put]
func (c *Publisher) UpdatePublisher(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Publisher
	if err := c.findByID(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.PublisherUpdatePayload
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
	render.JSON(w, r, resources.NewPublisherResource(item))
}

// DeletePublisher godoc
// @Summary Delete a Publisher
// @Tags Publisher
// @Security BearerAuth
// @Param id path int true "Publisher ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /publishers/{id} [delete]
func (c *Publisher) DeletePublisher(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Publisher
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

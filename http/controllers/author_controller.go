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

type Author struct {
	DB *gorm.DB
}

func (c *Author) findByID(id uint, item *models.Author) error {
	return c.DB.First(item, id).Error
}

// GetAllAuthor godoc
// @Summary List all Author records
// @Description Get a paginated list of Author records
// @Tags Author
// @Produce json
// @Security BearerAuth
// @Param page     query int    false "Page number (default 1)"
// @Param per_page query int    false "Items per page (default 15)"
// @Param search   query string false "Search by name"
// @Success 200 {object} resources.AuthorCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /authors [get]
func (c *Author) GetAllAuthor(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	search := r.URL.Query().Get("search")

	q := resources.NewAuthorQuery(c.DB)
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

// GetAuthorByID godoc
// @Summary Get an Author by ID
// @Tags Author
// @Produce json
// @Security BearerAuth
// @Param id path int true "Author ID"
// @Success 200 {object} resources.AuthorResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /authors/{id} [get]
func (c *Author) GetAuthorByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewAuthorQuery(c.DB).Find(id)
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

// CreateAuthor godoc
// @Summary Create an Author
// @Tags Author
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.AuthorPayload true "Author payload"
// @Success 201 {object} resources.AuthorResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /authors [post]
func (c *Author) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var payload requests.AuthorPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	item := models.Author{Name: payload.Name}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resources.NewAuthorResource(item))
}

// UpdateAuthor godoc
// @Summary Update an Author
// @Tags Author
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Author ID"
// @Param payload body requests.AuthorUpdatePayload true "Author update payload"
// @Success 200 {object} resources.AuthorResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /authors/{id} [put]
func (c *Author) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Author
	if err := c.findByID(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.AuthorUpdatePayload
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
	render.JSON(w, r, resources.NewAuthorResource(item))
}

// DeleteAuthor godoc
// @Summary Delete an Author
// @Tags Author
// @Security BearerAuth
// @Param id path int true "Author ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /authors/{id} [delete]
func (c *Author) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Author
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

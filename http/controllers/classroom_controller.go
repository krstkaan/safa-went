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

type Classroom struct {
	DB *gorm.DB
}

func (c *Classroom) findByID(id uint, item *models.Classroom) error {
	return c.DB.First(item, id).Error
}

// GetAllClassroom godoc
// @Summary List all Classroom records
// @Description Get a paginated list of Classroom records
// @Tags Classroom
// @Produce json
// @Security BearerAuth
// @Param page     query int    false "Page number (default 1)"
// @Param per_page query int    false "Items per page (default 15)"
// @Param search   query string false "Search by name"
// @Success 200 {object} resources.ClassroomCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /classrooms [get]
func (c *Classroom) GetAllClassroom(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	search := r.URL.Query().Get("search")

	q := resources.NewClassroomQuery(c.DB)
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

// GetClassroomByID godoc
// @Summary Get a Classroom by ID
// @Tags Classroom
// @Produce json
// @Security BearerAuth
// @Param id path int true "Classroom ID"
// @Success 200 {object} resources.ClassroomResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /classrooms/{id} [get]
func (c *Classroom) GetClassroomByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewClassroomQuery(c.DB).Find(id)
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

// CreateClassroom godoc
// @Summary Create a Classroom
// @Tags Classroom
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.ClassroomPayload true "Classroom payload"
// @Success 201 {object} resources.ClassroomResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /classrooms [post]
func (c *Classroom) CreateClassroom(w http.ResponseWriter, r *http.Request) {
	var payload requests.ClassroomPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	item := models.Classroom{Name: payload.Name}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resources.NewClassroomResource(item))
}

// UpdateClassroom godoc
// @Summary Update a Classroom
// @Tags Classroom
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Classroom ID"
// @Param payload body requests.ClassroomUpdatePayload true "Classroom update payload"
// @Success 200 {object} resources.ClassroomResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /classrooms/{id} [put]
func (c *Classroom) UpdateClassroom(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Classroom
	if err := c.findByID(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.ClassroomUpdatePayload
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
	render.JSON(w, r, resources.NewClassroomResource(item))
}

// DeleteClassroom godoc
// @Summary Delete a Classroom
// @Tags Classroom
// @Security BearerAuth
// @Param id path int true "Classroom ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /classrooms/{id} [delete]
func (c *Classroom) DeleteClassroom(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Classroom
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

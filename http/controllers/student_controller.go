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

type Student struct {
	DB *gorm.DB
}

func (c *Student) findByID(id uint, item *models.Student) error {
	return c.DB.Preload("Classroom").First(item, id).Error
}

// GetAllStudent godoc
// @Summary List all Student records
// @Description Get a paginated list of Student records with optional filters
// @Tags Student
// @Produce json
// @Security BearerAuth
// @Param page         query int    false "Page number (default 1)"
// @Param per_page     query int    false "Items per page (default 15)"
// @Param search       query string false "Search by name"
// @Param classroom_id query int    false "Filter by classroom ID"
// @Success 200 {object} resources.StudentCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /students [get]
func (c *Student) GetAllStudent(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	search := r.URL.Query().Get("search")
	classroomIDStr := r.URL.Query().Get("classroom_id")
	var classroomID uint
	if v, err := strconv.ParseUint(classroomIDStr, 10, 64); err == nil {
		classroomID = uint(v)
	}

	q := resources.NewStudentQuery(c.DB)
	collection, err := q.Filter(search, classroomID, page, perPage)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, collection)
}

// GetStudentByID godoc
// @Summary Get a Student by ID
// @Tags Student
// @Produce json
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Success 200 {object} resources.StudentResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /students/{id} [get]
func (c *Student) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewStudentQuery(c.DB).Find(id)
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

// CreateStudent godoc
// @Summary Create a Student
// @Tags Student
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.StudentPayload true "Student payload"
// @Success 201 {object} resources.StudentResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /students [post]
func (c *Student) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var payload requests.StudentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	item := models.Student{Name: payload.Name, ClassroomID: payload.ClassroomID}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	created, err := resources.NewStudentQuery(c.DB).Find(item.ID)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, created)
}

// UpdateStudent godoc
// @Summary Update a Student
// @Tags Student
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Param payload body requests.StudentUpdatePayload true "Student update payload"
// @Success 200 {object} resources.StudentResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /students/{id} [put]
func (c *Student) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Student
	if err := c.findByID(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.StudentUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	updates := map[string]interface{}{}
	if payload.Name != nil {
		updates["name"] = *payload.Name
	}
	if payload.ClassroomID != nil {
		updates["classroom_id"] = *payload.ClassroomID
	}
	if len(updates) == 0 {
		responses.JSONError(w, r, http.StatusBadRequest, "no fields to update")
		return
	}
	if err := c.DB.Model(&item).Updates(updates).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	updated, err := resources.NewStudentQuery(c.DB).Find(id)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, updated)
}

// DeleteStudent godoc
// @Summary Delete a Student
// @Tags Student
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /students/{id} [delete]
func (c *Student) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Student
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

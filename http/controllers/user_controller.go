package controllers

import (
	"encoding/json"
	"net/http"
	"safa-went/database/models"
	"safa-went/http/requests"
	"safa-went/http/resources"
	"safa-went/internal/responses"
	"strconv"

	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func (c *User) findByID(id uint, item *models.User) error {
	return c.DB.First(item, id).Error
}

// GetAllUser godoc
// @Summary List all User records
// @Description Get a paginated list of User records. Use ?page=1&per_page=15 to control pagination.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param page     query int false "Page number (default 1)"
// @Param per_page query int false "Items per page (default 15)"
// @Success 200 {object} resources.UserCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /user [get]
func (c *User) GetAllUser(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)

	collection, err := resources.NewUserQuery(c.DB).Paginate(page, perPage)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, collection)
}

// GetUserByID godoc
// @Summary Get a User by ID
// @Description Get a single User record by ID
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} resources.UserResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /user/{id} [get]
func (c *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := resources.NewUserQuery(c.DB).Find(id)
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

// CreateUser godoc
// @Summary Create a User
// @Description Create a new User record
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.UserPayload true "User payload"
// @Success 201 {object} resources.UserResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /user [post]
func (c *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload requests.UserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	item := models.User{
		Name:  payload.Name,
		Email: payload.Email,
	}

	err = c.DB.Create(&item).Error
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resources.NewUserResource(item))
}

// UpdateUser godoc
// @Summary Update a User
// @Description Update an existing User record by ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param payload body requests.UserUpdatePayload true "User update payload"
// @Success 200 {object} resources.UserResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /user/{id} [put]
func (c *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	var item models.User
	err = c.findByID(id, &item)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	var payload requests.UserUpdatePayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	updates := map[string]interface{}{}
	if payload.Name != nil {
		updates["name"] = *payload.Name
	}
	if payload.Email != nil {
		updates["email"] = *payload.Email
	}
	if len(updates) == 0 {
		responses.JSONError(w, r, http.StatusBadRequest, "no fields to update")
		return
	}

	err = c.DB.Model(&item).Updates(updates).Error
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = c.findByID(id, &item)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, resources.NewUserResource(item))
}

// DeleteUser godoc
// @Summary Delete a User
// @Description Delete an existing User record by ID
// @Tags User
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /user/{id} [delete]
func (c *User) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	err = c.DB.Delete(&models.User{}, id).Error
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusNoContent)
}

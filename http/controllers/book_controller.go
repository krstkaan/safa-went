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

type Book struct {
	DB *gorm.DB
}

func (c *Book) findRaw(id uint, item *models.Book) error {
	return c.DB.Preload("Author").Preload("Publisher").First(item, id).Error
}

// GetAllBook godoc
// @Summary List all Book records
// @Description Get a paginated list of books with optional filters. Default sort: fixture_no ASC.
// @Tags Book
// @Produce json
// @Security BearerAuth
// @Param page         query int    false "Page number (default 1)"
// @Param per_page     query int    false "Items per page (default 15)"
// @Param name         query string false "Filter by name (partial match)"
// @Param fixture_no   query int    false "Filter by fixture number (exact)"
// @Param search       query string false "Search by name or barcode"
// @Param author_id    query int    false "Filter by author ID"
// @Param publisher_id query int    false "Filter by publisher ID"
// @Param level        query string false "Filter by level (ilkokul/ortaokul/ortak)"
// @Param is_donation  query bool   false "Filter by donation flag"
// @Param sort         query string false "Sort column (e.g. fixture_no ASC)"
// @Success 200 {object} resources.BookCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /books [get]
func (c *Book) GetAllBook(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)

	name := r.URL.Query().Get("name")
	fixtureNoStr := r.URL.Query().Get("fixture_no")
	var fixtureNo int
	if v, err := strconv.Atoi(fixtureNoStr); err == nil {
		fixtureNo = v
	}
	search := r.URL.Query().Get("search")
	authorIDStr := r.URL.Query().Get("author_id")
	var authorID uint
	if v, err := strconv.ParseUint(authorIDStr, 10, 64); err == nil {
		authorID = uint(v)
	}
	publisherIDStr := r.URL.Query().Get("publisher_id")
	var publisherID uint
	if v, err := strconv.ParseUint(publisherIDStr, 10, 64); err == nil {
		publisherID = uint(v)
	}
	level := r.URL.Query().Get("level")
	sortBy := r.URL.Query().Get("sort")

	var isDonation *bool
	if v := r.URL.Query().Get("is_donation"); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			isDonation = &b
		}
	}

	q := resources.NewBookQuery(c.DB)
	collection, err := q.Filter(name, fixtureNo, search, authorID, publisherID, level, isDonation, sortBy, page, perPage)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, collection)
}

// GetBookByID godoc
// @Summary Get a Book by ID
// @Tags Book
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Success 200 {object} resources.BookResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /books/{id} [get]
func (c *Book) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewBookQuery(c.DB).Find(id)
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

// CreateBook godoc
// @Summary Create a Book
// @Tags Book
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.BookPayload true "Book payload"
// @Success 201 {object} resources.BookResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /books [post]
func (c *Book) CreateBook(w http.ResponseWriter, r *http.Request) {
	var payload requests.BookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	item := models.Book{
		Name:        payload.Name,
		Barcode:     payload.Barcode,
		AuthorID:    payload.AuthorID,
		PublisherID: payload.PublisherID,
		Language:    payload.Language,
		PageCount:   payload.PageCount,
		IsDonation:  payload.IsDonation,
		ShelfCode:   payload.ShelfCode,
		FixtureNo:   payload.FixtureNo,
		Level:       payload.Level,
	}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	created, err := resources.NewBookQuery(c.DB).Find(item.ID)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, created)
}

// UpdateBook godoc
// @Summary Update a Book
// @Tags Book
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Param payload body requests.BookUpdatePayload true "Book update payload"
// @Success 200 {object} resources.BookResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /books/{id} [put]
func (c *Book) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Book
	if err := c.findRaw(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.BookUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	updates := map[string]interface{}{}
	if payload.Name != nil {
		updates["name"] = *payload.Name
	}
	if payload.Barcode != nil {
		updates["barcode"] = *payload.Barcode
	}
	if payload.AuthorID != nil {
		updates["author_id"] = *payload.AuthorID
	}
	if payload.PublisherID != nil {
		updates["publisher_id"] = *payload.PublisherID
	}
	if payload.Language != nil {
		updates["language"] = *payload.Language
	}
	if payload.PageCount != nil {
		updates["page_count"] = *payload.PageCount
	}
	if payload.IsDonation != nil {
		updates["is_donation"] = *payload.IsDonation
	}
	if payload.ShelfCode != nil {
		updates["shelf_code"] = *payload.ShelfCode
	}
	if payload.FixtureNo != nil {
		updates["fixture_no"] = *payload.FixtureNo
	}
	if payload.Level != nil {
		updates["level"] = *payload.Level
	}
	if len(updates) == 0 {
		responses.JSONError(w, r, http.StatusBadRequest, "no fields to update")
		return
	}
	if err := c.DB.Model(&item).Updates(updates).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	updated, err := resources.NewBookQuery(c.DB).Find(id)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, updated)
}

// DeleteBook godoc
// @Summary Delete a Book
// @Tags Book
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /books/{id} [delete]
func (c *Book) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Book
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

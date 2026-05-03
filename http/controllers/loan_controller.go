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

type Loan struct {
	DB *gorm.DB
}

func (c *Loan) findRaw(id uint, item *models.Loan) error {
	return c.DB.
		Preload("Student").
		Preload("Student.Classroom").
		Preload("Book").
		Preload("Book.Author").
		Preload("Book.Publisher").
		First(item, id).Error
}

// GetAllLoan godoc
// @Summary List all Loan records
// @Description Get a paginated list of loans with optional filters
// @Tags Loan
// @Produce json
// @Security BearerAuth
// @Param page       query int    false "Page number (default 1)"
// @Param per_page   query int    false "Items per page (default 15)"
// @Param status     query string false "Filter by status (active/returned/overdue)"
// @Param search     query string false "Search by student name"
// @Param student_id query int    false "Filter by student ID"
// @Param book_id    query int    false "Filter by book ID"
// @Param sort       query string false "Sort column"
// @Success 200 {object} resources.LoanCollection
// @Failure 500 {object} responses.ErrorBody
// @Router /loans [get]
func (c *Loan) GetAllLoan(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort")

	var studentID, bookID uint
	if v, err := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64); err == nil {
		studentID = uint(v)
	}
	if v, err := strconv.ParseUint(r.URL.Query().Get("book_id"), 10, 64); err == nil {
		bookID = uint(v)
	}

	q := resources.NewLoanQuery(c.DB)
	collection, err := q.Filter(status, search, studentID, bookID, sortBy, page, perPage)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, collection)
}

// GetLoanByID godoc
// @Summary Get a Loan by ID
// @Tags Loan
// @Produce json
// @Security BearerAuth
// @Param id path int true "Loan ID"
// @Success 200 {object} resources.LoanResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /loans/{id} [get]
func (c *Loan) GetLoanByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := resources.NewLoanQuery(c.DB).Find(id)
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

// CreateLoan godoc
// @Summary Create a Loan
// @Description Creates a new loan after checking level compatibility and active-loan constraints.
// @Tags Loan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.LoanPayload true "Loan payload"
// @Success 201 {object} resources.LoanResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 422 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /loans [post]
func (c *Loan) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var payload requests.LoanPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	// Load student with classroom
	var student models.Student
	if err := c.DB.Preload("Classroom").First(&student, payload.StudentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusBadRequest, "student not found")
		} else {
			responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Load book
	var book models.Book
	if err := c.DB.First(&book, payload.BookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusBadRequest, "book not found")
		} else {
			responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Check level compatibility
	if eligible, reason := models.CanStudentBorrowBook(student, book); !eligible {
		responses.JSONError(w, r, http.StatusUnprocessableEntity, reason)
		return
	}

	// Check for active loan on this book
	var activeCount int64
	c.DB.Model(&models.Loan{}).
		Where("book_id = ? AND status = ? AND deleted_at IS NULL", payload.BookID, models.LoanStatusActive).
		Count(&activeCount)
	if activeCount > 0 {
		responses.JSONError(w, r, http.StatusUnprocessableEntity, "book is currently loaned out")
		return
	}

	item := models.Loan{
		StudentID: payload.StudentID,
		BookID:    payload.BookID,
		LoanDate:  payload.LoanDate,
		DueDate:   payload.DueDate,
		Status:    models.LoanStatusActive,
		Notes:     payload.Notes,
	}
	if err := c.DB.Create(&item).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	created, err := resources.NewLoanQuery(c.DB).Find(item.ID)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, created)
}

// UpdateLoan godoc
// @Summary Update a Loan
// @Tags Loan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Loan ID"
// @Param payload body requests.LoanUpdatePayload true "Loan update payload"
// @Success 200 {object} resources.LoanResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /loans/{id} [put]
func (c *Loan) UpdateLoan(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Loan
	if err := c.findRaw(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	var payload requests.LoanUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	updates := map[string]interface{}{}
	if payload.DueDate != nil {
		updates["due_date"] = *payload.DueDate
	}
	if payload.Status != nil {
		updates["status"] = *payload.Status
	}
	if payload.Notes != nil {
		updates["notes"] = *payload.Notes
	}
	if len(updates) == 0 {
		responses.JSONError(w, r, http.StatusBadRequest, "no fields to update")
		return
	}
	if err := c.DB.Model(&item).Updates(updates).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	updated, err := resources.NewLoanQuery(c.DB).Find(id)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, updated)
}

// DeleteLoan godoc
// @Summary Delete a Loan
// @Tags Loan
// @Security BearerAuth
// @Param id path int true "Loan ID"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Router /loans/{id} [delete]
func (c *Loan) DeleteLoan(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Loan
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

// ReturnBook godoc
// @Summary Mark a loan as returned
// @Tags Loan
// @Produce json
// @Security BearerAuth
// @Param id path int true "Loan ID"
// @Success 200 {object} resources.LoanResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Failure 422 {object} responses.ErrorBody
// @Router /loans/{id}/return [post]
func (c *Loan) ReturnBook(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	var item models.Loan
	if err := c.findRaw(id, &item); err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if item.Status == models.LoanStatusReturned {
		responses.JSONError(w, r, http.StatusUnprocessableEntity, "loan is already returned")
		return
	}
	now := time.Now()
	if err := c.DB.Model(&item).Updates(map[string]interface{}{
		"status":      models.LoanStatusReturned,
		"return_date": now,
	}).Error; err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	updated, err := resources.NewLoanQuery(c.DB).Find(id)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, updated)
}

// CheckAvailability godoc
// @Summary Check if a student can borrow a book
// @Tags Loan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.LoanCheckAvailabilityPayload true "Check payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} responses.ErrorBody
// @Router /loans/check-availability [post]
func (c *Loan) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var payload requests.LoanCheckAvailabilityPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	var student models.Student
	if err := c.DB.Preload("Classroom").First(&student, payload.StudentID).Error; err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "student not found")
		return
	}

	var book models.Book
	if err := c.DB.First(&book, payload.BookID).Error; err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "book not found")
		return
	}

	eligible, reason := models.CanStudentBorrowBook(student, book)
	if !eligible {
		render.JSON(w, r, map[string]interface{}{"available": false, "reason": reason})
		return
	}

	var activeCount int64
	c.DB.Model(&models.Loan{}).
		Where("book_id = ? AND status = ? AND deleted_at IS NULL", payload.BookID, models.LoanStatusActive).
		Count(&activeCount)
	if activeCount > 0 {
		render.JSON(w, r, map[string]interface{}{"available": false, "reason": "book is currently loaned out"})
		return
	}

	render.JSON(w, r, map[string]interface{}{"available": true, "reason": ""})
}

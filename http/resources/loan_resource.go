package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type LoanResource struct {
	ID         uint            `json:"id"`
	StudentID  uint            `json:"student_id"`
	Student    StudentResource `json:"student"`
	BookID     uint            `json:"book_id"`
	Book       BookResource    `json:"book"`
	LoanDate   time.Time       `json:"loan_date"`
	DueDate    time.Time       `json:"due_date"`
	ReturnDate *time.Time      `json:"return_date"`
	Status     string          `json:"status"`
	Notes      *string         `json:"notes"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type LoanCollection struct {
	Data []LoanResource `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type LoanQuery struct {
	db *gorm.DB
}

func NewLoanQuery(db *gorm.DB) *LoanQuery {
	return &LoanQuery{db: db}
}

func (q *LoanQuery) newLoanResource(m models.Loan) LoanResource {
	return LoanResource{
		ID:         m.ID,
		StudentID:  m.StudentID,
		Student:    NewStudentResource(m.Student),
		BookID:     m.BookID,
		Book:       newBookResourceWithLoan(m.Book, q.db),
		LoanDate:   m.LoanDate,
		DueDate:    m.DueDate,
		ReturnDate: m.ReturnDate,
		Status:     m.Status,
		Notes:      m.Notes,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func (q *LoanQuery) preloaded() *gorm.DB {
	return q.db.
		Preload("Student").
		Preload("Student.Classroom").
		Preload("Book").
		Preload("Book.Author").
		Preload("Book.Publisher")
}

func (q *LoanQuery) Find(id uint) (LoanResource, error) {
	var item models.Loan
	if err := q.preloaded().First(&item, id).Error; err != nil {
		return LoanResource{}, err
	}
	return q.newLoanResource(item), nil
}

// Filter applies optional filters and returns a paginated collection.
func (q *LoanQuery) Filter(status string, search string, studentID, bookID uint, sortBy string, page, perPage int64) (LoanCollection, error) {
	db := q.preloaded()
	if status != "" {
		db = db.Where("loans.status = ?", status)
	}
	if search != "" {
		db = db.Joins("JOIN students ON students.id = loans.student_id").
			Where("students.name ILIKE ?", "%"+search+"%")
	}
	if studentID > 0 {
		db = db.Where("loans.student_id = ?", studentID)
	}
	if bookID > 0 {
		db = db.Where("loans.book_id = ?", bookID)
	}
	if sortBy != "" {
		db = db.Order(sortBy)
	} else {
		db = db.Order("loans.loan_date DESC")
	}

	var items []models.Loan
	var total int64
	countDB := db.Session(&gorm.Session{})
	if err := countDB.Model(&models.Loan{}).Count(&total).Error; err != nil {
		return LoanCollection{}, err
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage
	if err := db.Offset(int(offset)).Limit(int(perPage)).Find(&items).Error; err != nil {
		return LoanCollection{}, err
	}

	data := make([]LoanResource, 0, len(items))
	for _, item := range items {
		data = append(data, q.newLoanResource(item))
	}

	meta := BuildMeta(total, page, perPage, int64(len(items)))
	return LoanCollection{Data: data, Meta: meta}, nil
}

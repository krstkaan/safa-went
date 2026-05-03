package models

import (
	"time"

	"gorm.io/gorm"
)

// LoanStatus enum values.
const (
	LoanStatusActive   = "active"
	LoanStatusReturned = "returned"
	LoanStatusOverdue  = "overdue"
)

type Loan struct {
	gorm.Model
	StudentID  uint       `json:"student_id" gorm:"not null"`
	Student    Student    `json:"student" gorm:"foreignKey:StudentID"`
	BookID     uint       `json:"book_id" gorm:"not null"`
	Book       Book       `json:"book" gorm:"foreignKey:BookID"`
	LoanDate   time.Time  `json:"loan_date" gorm:"not null"`
	DueDate    time.Time  `json:"due_date" gorm:"not null"`
	ReturnDate *time.Time `json:"return_date"`
	Status     string     `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
	Notes      *string    `json:"notes"`
}

// CanStudentBorrowBook checks business rules for borrowing eligibility.
// Returns (eligible, reason). Caller must ensure Student.Classroom is preloaded.
func CanStudentBorrowBook(student Student, book Book) (bool, string) {
	studentLevel := student.GetLevel()
	if studentLevel != book.Level && book.Level != BookLevelOrtak && studentLevel != BookLevelOrtak {
		return false, "book level does not match student level"
	}
	return true, ""
}

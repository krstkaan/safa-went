package requests

import "time"

type LoanPayload struct {
	StudentID uint      `json:"student_id"`
	BookID    uint      `json:"book_id"`
	LoanDate  time.Time `json:"loan_date"`
	DueDate   time.Time `json:"due_date"`
	Notes     *string   `json:"notes,omitempty"`
}

type LoanUpdatePayload struct {
	DueDate *time.Time `json:"due_date,omitempty"`
	Status  *string    `json:"status,omitempty"`
	Notes   *string    `json:"notes,omitempty"`
}

type LoanCheckAvailabilityPayload struct {
	StudentID uint `json:"student_id"`
	BookID    uint `json:"book_id"`
}

package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	ClassroomID uint      `json:"classroom_id" gorm:"not null"`
	Classroom   Classroom `json:"classroom" gorm:"foreignKey:ClassroomID"`
}

// GetLevel delegates to the associated Classroom.
func (s *Student) GetLevel() string {
	return s.Classroom.GetLevel()
}
